package repo

import (
	"context"
	"github.com/Zkeai/go_template/internal/conf"
	"github.com/Zkeai/go_template/internal/repo/db"
)

type Repo struct {
	db *db.DB
}

func NewRepo(conf *conf.Conf) *Repo {
	return &Repo{
		db: db.NewDB(conf.DB),
	}
}

func (r *Repo) UserRegister(ctx context.Context, wallet string) (*db.YuUser, error) {

	return r.db.InsertUser(ctx, wallet)
}

func (r *Repo) UserQuery(ctx context.Context, wallet string) (*db.YuUser, error) {

	return r.db.QueryUser(ctx, wallet)
}
