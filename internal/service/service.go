package service

import "github.com/Zkeai/go_template/internal/conf"
import "github.com/Zkeai/go_template/internal/repo"

type Service struct {
	conf *conf.Conf
	repo *repo.Repo
}

func NewService(conf *conf.Conf) *Service {
	return &Service{
		conf: conf,
		repo: repo.NewRepo(conf),
	}
}
