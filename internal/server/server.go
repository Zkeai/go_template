package server

import (
	chttp "github.com/Zkeai/go_template/common/net/cttp"
	"github.com/Zkeai/go_template/internal/conf"

	"github.com/Zkeai/go_template/internal/handler"
	"github.com/Zkeai/go_template/internal/service"
)

func NewHTTP(conf *conf.Conf) *chttp.Server {
	s := chttp.NewServer(conf.Server)

	svc := service.NewService(conf)
	handler.InitRouter(s, svc)

	err := s.Start()

	if err != nil {
		panic(err)
	}

	return s
}
