package main

import (
	"context"
	"flag"
	cconf "github.com/Zkeai/go_template/common/conf"
	"github.com/Zkeai/go_template/common/cron"
	"github.com/Zkeai/go_template/common/logger"
	"github.com/Zkeai/go_template/common/redis"
	"github.com/Zkeai/go_template/common/util"
	"github.com/Zkeai/go_template/internal/conf"
	"github.com/Zkeai/go_template/internal/server"
	"github.com/ouqiang/goutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	// filePath yaml文件目录
	filePath *string
	// AppDir 应用根目录
	AppDir string
	// LogDir 日志目录
	LogDir string // 日志目录

)

// @title		MuPay
// @version		1.0.0
// @description	木鱼发卡 https://github.com/zkeai
// @host			localhost:2900
// @BasePath		/api/v1

// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization
func main() {
	//初始化配置
	initEnv()

}

func initEnv() {
	//logger 初始化
	logger.InitLogger()
	AppDir, err := goutil.WorkDir()
	if err != nil {
		logger.Fatal(err)
	}
	LogDir = filepath.Join(AppDir, "/log")
	util.CreateDirIfNotExists(LogDir)

	//读取yaml配置
	flag.Parse()
	filePath = flag.String("conf", "etc/config.yaml", "the config path")
	c := new(conf.Conf)
	err = cconf.Unmarshal(*filePath, c)
	if err != nil {
		logger.Error(err)
	}

	//redis 初始化
	redis.InitRedis(c.Rdb)

	//cron 初始化
	taskManager := cron.GetManager()
	defer taskManager.Stop() // 确保在程序退出前停止调度器

	//http 初始化
	srv := server.NewHTTP(c)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			_ = srv.Shutdown(context.Background())
			return
		default:

			return
		}
	}
}
