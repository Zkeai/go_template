package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Zkeai/go_template/common/logger"
	"time"
)

// YuUser 用户表
type YuUser struct {
	ID         int64      `json:"id" gorm:"id"`                   // id
	Wallet     string     `json:"wallet" gorm:"wallet"`           // eth 钱包地址
	Status     int8       `json:"status" gorm:"status"`           // 账号状态
	Type       int64      `json:"type" gorm:"type"`               // 账号类型 0-用户 1-商户 2-管理员
	CreateTime time.Time  `json:"create_time" gorm:"create_time"` // 创建时间
	LoginTime  *time.Time `json:"login_time" gorm:"login_time"`   // 登录时间
	LoginIp    *string    `json:"login_ip" gorm:"login_ip"`       // 登录ip
}

const (
	insertManageSQL = `INSERT INTO yu_user (wallet) VALUES (?)`
	queryManageSql  = `SELECT * FROM yu_user WHERE wallet = ?`
)

func (db *DB) InsertUser(ctx context.Context, wallet string) (*YuUser, error) {

	_, err := db.db.Exec(ctx, insertManageSQL, wallet)
	if err != nil {
		return nil, err
	}
	user, err := db.QueryUser(ctx, wallet)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) QueryUser(ctx context.Context, wallet string) (*YuUser, error) {
	row := db.db.QueryRow(ctx, queryManageSql, wallet)
	u := &YuUser{}
	err := row.Scan(&u.ID, &u.Wallet, &u.Status, &u.Type, &u.CreateTime, &u.LoginTime, &u.LoginIp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.Error(err)
		return nil, err
	}
	return u, err
}
