package handler

import (
	"github.com/Zkeai/go_template/common/conf"
	"github.com/Zkeai/go_template/common/middleware"
	chttp "github.com/Zkeai/go_template/common/net/cttp"
	_ "github.com/Zkeai/go_template/docs"
	"github.com/Zkeai/go_template/internal/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"

	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

var svc *service.Service

func InitRouter(s *chttp.Server, service *service.Service) {
	svc = service
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// 初始化 Casbin Enforcer
	e, err := casbin.NewEnforcer(dir+"/common/conf/rbac_model.conf", dir+"/common/conf/rabc_policy.csv")
	if err != nil {
		log.Fatalf("Failed to initialize Casbin enforcer: %v", err)
	}

	g := s.Group("/api/v1")
	g.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	g.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, conf.Response{Code: 200, Msg: "success", Data: "木鱼喜欢你"})
	})

	ug := g.Group("/user")
	ugpub := ug.Group("/public")
	{
		ugpub.POST("/register", userRegister)
		ugpub.POST("/login", userLogin)
	}
	ugpro := ug.Group("/protected")

	ugpro.Use(middleware.Middleware())
	ugpro.Use(middleware.CasbinMiddleware(e))
	{
		ugpro.GET("/query", userQuery)
	}

}
