package main

import (
	"fmt"

	"test.com/project-user/config"

	"github.com/gin-gonic/gin"
	"test.com/common"
	"test.com/common/logs"
	_ "test.com/project-user/api/user"
	"test.com/project-user/router"
)

func main() {
	r := gin.New()
	r.Use(logs.GinLogger(), logs.GinRecovery(true))

	//路由
	router.InitRouter(r)
	//grpc服务注册
	grpcServer := router.RegisterGrpc()
	common.Run(r, config.AppConf.App.Name, fmt.Sprintf(":%d", config.AppConf.App.Port))
}
