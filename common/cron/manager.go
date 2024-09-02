package cron

import (
	"context"
	"fmt"
	"github.com/Zkeai/go_template/common/logger"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
)

type TaskManager struct {
	c           *cron.Cron
	once        sync.Once
	mu          sync.Mutex // 保护tasks和onceTasks映射
	tasks       map[cron.EntryID]context.CancelFunc
	onceTasks   map[int]context.CancelFunc
	WsConn      map[string]*websocket.Conn    // 用于存储每个用户的 WebSocket 连接
	cancelFuncs map[string]context.CancelFunc // 用于存储每个用户的 context.CancelFunc
	Ctx         context.Context
}

var manager *TaskManager

// GetManager 获取单例的TaskManager
func GetManager() *TaskManager {
	if manager == nil {
		manager = &TaskManager{
			tasks:       make(map[cron.EntryID]context.CancelFunc),
			onceTasks:   make(map[int]context.CancelFunc),
			WsConn:      make(map[string]*websocket.Conn),
			cancelFuncs: make(map[string]context.CancelFunc),
		}
		manager.once.Do(func() {
			manager.c = cron.New(cron.WithSeconds())
			manager.c.Start()
		})
	}
	return manager
}

// AddTask 注册新的定时任务
func (tm *TaskManager) AddTask(spec string, cmd func(context.Context)) (cron.EntryID, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	entryID, err := tm.c.AddFunc(spec, func() {
		select {
		case <-ctx.Done():
			logger.Fatal("Task cancelled:", ctx.Err())
		default:
			cmd(ctx)
		}
	})
	if err != nil {
		cancel() // 如果添加任务失败，则取消context
		logger.Error("Task add err", err)
		return 0, err
	}
	tm.tasks[entryID] = cancel
	logger.Info("Task add success:", entryID)
	return entryID, nil
}

// AddTaskOnce 注册只运行一次的任务
func (tm *TaskManager) AddTaskOnce(d time.Duration, cmd func(context.Context)) int {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	taskID := len(tm.onceTasks) + 1
	tm.onceTasks[taskID] = cancel

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled:", ctx.Err())
		case <-time.After(d):
			cmd(ctx)
			tm.mu.Lock()
			delete(tm.onceTasks, taskID)
			tm.mu.Unlock()
		}
	}()

	return taskID
}

// ConnectWebSocket 连接 WebSocket 并保存连接
func (tm *TaskManager) ConnectWebSocket(userID, url string) (*TaskManager, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	tm.WsConn[userID] = conn

	// 启动一个 goroutine 来读取消息并处理
	ctx, cancel := context.WithCancel(context.Background())
	tm.cancelFuncs[userID] = cancel
	tm.Ctx = ctx

	return tm, nil
}

// CloseWebSocket 关闭指定用户的 WebSocket 连接
func (tm *TaskManager) CloseWebSocket(userID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if conn, ok := tm.WsConn[userID]; ok {
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			logger.Error("CloseWebSocket close:", err)
			return
		}
		conn.Close()
		delete(tm.WsConn, userID)
	}
	if cancel, ok := tm.cancelFuncs[userID]; ok {
		cancel()
		delete(tm.cancelFuncs, userID)
	}
}

// RemoveTask 删除定时任务
func (tm *TaskManager) RemoveTask(entryID cron.EntryID) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if cancel, ok := tm.tasks[entryID]; ok {
		cancel()             // 取消任务
		tm.c.Remove(entryID) // 从cron调度器中移除任务
		delete(tm.tasks, entryID)
		logger.Info("Task removed:", entryID)
	}
}

// Stop 停止任务管理器
func (tm *TaskManager) Stop() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for entryID, cancel := range tm.tasks {
		cancel() // 取消所有任务
		tm.c.Remove(entryID)
	}
	tm.tasks = make(map[cron.EntryID]context.CancelFunc)
	tm.c.Stop()
	for userID := range tm.WsConn {
		tm.CloseWebSocket(userID)
	}
	logger.Info("Cron stopped")
}
