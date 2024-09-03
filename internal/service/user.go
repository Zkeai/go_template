package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Zkeai/go_template/common/logger"
	"github.com/Zkeai/go_template/common/middleware"
	"github.com/Zkeai/go_template/common/redis"
	"github.com/Zkeai/go_template/internal/repo/db"
	redisv8 "github.com/go-redis/redis/v8"

	"io"
	"time"
)

type UserRegisterResponse struct {
	UserExists bool       // 标识钱包是否重复
	User       *db.YuUser // 注册的用户信息
}

type SessionData struct {
	Role      int    `json:"role"`
	Status    int    `json:"status"`
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
}
type LoginResponse struct {
	SessionID string `json:"session_id"`
	Token     string `json:"token"`
}

func (s *Service) UserRegister(ctx context.Context, wallet string) (*UserRegisterResponse, error) {
	userModel, err := s.repo.UserQuery(ctx, wallet)
	if userModel != nil {
		// 返回 UserExists 为 true，并不返回错误
		return &UserRegisterResponse{
			UserExists: true,
			User:       userModel,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	user, err := s.repo.UserRegister(ctx, wallet)
	if err != nil {
		return nil, err
	}

	// 返回正常的注册结果，UserExists 为 false
	return &UserRegisterResponse{
		UserExists: false,
		User:       user,
	}, nil
}

func (s *Service) UserLogin(ctx context.Context, wallet string) (*LoginResponse, error) {

	query, err := s.repo.UserQuery(ctx, wallet)
	if err != nil || query == nil {
		return &LoginResponse{
			SessionID: "",
			Token:     "用户不存在",
		}, err
	}
	//生成sessionID
	id, err := generateSessionID()
	if err != nil {
		return &LoginResponse{}, err
	}

	//生成jwt
	token, err := middleware.GenerateToken(wallet, id)
	if err != nil {
		return &LoginResponse{}, err
	}
	//redis取jwt
	result, err := redis.GetClient().Get(ctx, wallet).Result()

	var sessionData SessionData
	if errors.Is(err, redisv8.Nil) {
	} else {
		err = json.Unmarshal([]byte(result), &sessionData)
		if err != nil {
			logger.Error("Failed to unmarshal JSON data: %v", err)
			return nil, fmt.Errorf("failed to unmarshal")
		}
		_ = middleware.InvalidateToken(sessionData.Token)
	}

	//redis存数据
	userData := SessionData{
		Role:      int(query.Type),
		Status:    int(query.Status),
		Token:     token,
		SessionID: id,
	}
	// 将数据转换为 JSON 字符串
	jsonData, err := json.Marshal(userData)
	if err != nil {
		return &LoginResponse{}, err
	}

	_, err = redis.GetClient().Set(ctx, query.Wallet, jsonData, time.Minute*10).Result()
	if err != nil {
		return &LoginResponse{}, err
	}
	return &LoginResponse{
		SessionID: id,
		Token:     token,
	}, nil
}

func (s *Service) UserQuery(ctx context.Context, wallet string) (*db.YuUser, error) {

	userModel, err := s.repo.UserQuery(ctx, wallet)
	if err != nil {
		return nil, err
	}
	return userModel, nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil // 使用Base64编码生成唯一的Session ID
}
