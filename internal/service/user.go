package service

import (
	"context"
	"encoding/json"
	"github.com/Zkeai/go_template/common/middleware"
	"github.com/Zkeai/go_template/common/redis"
	"github.com/Zkeai/go_template/internal/repo/db"
	"time"
)

type UserRegisterResponse struct {
	UserExists bool       // 标识钱包是否重复
	User       *db.YuUser // 注册的用户信息
}

type UserData struct {
	Role   int    `json:"role"`
	Status int    `json:"status"`
	Token  string `json:"token"`
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

func (s *Service) UserLogin(ctx context.Context, wallet string) (string, error) {

	query, err := s.repo.UserQuery(ctx, wallet)
	if err != nil || query == nil {
		return "用户不存在", err
	}

	//生成jwt
	token, err := middleware.GenerateToken(wallet)
	if err != nil {
		return "", err
	}
	//redis存数据
	userData := UserData{
		Role:   int(query.Type),
		Status: int(query.Status),
		Token:  token,
	}
	// 将数据转换为 JSON 字符串
	jsonData, err := json.Marshal(userData)
	if err != nil {
		return "", err
	}

	_, err = redis.GetClient().Set(ctx, query.Wallet, jsonData, time.Hour*24).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) UserQuery(ctx context.Context, wallet string) (*db.YuUser, error) {

	userModel, err := s.repo.UserQuery(ctx, wallet)
	if err != nil {
		return nil, err
	}
	return userModel, nil
}
