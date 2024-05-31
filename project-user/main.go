package main

import (
	"test.com/project-user/config"

	"github.com/gin-gonic/gin"
	"test.com/common"
	"test.com/common/logs"
	_ "test.com/project-user/internal/dao"
	"test.com/project-user/router"
)

func main() {
	r := gin.New()
	r.Use(logs.GinLogger(), logs.GinRecovery(true))
	//config.AppConf.InitZapLog()
	//zap.L().Info("初始化成功")
	//路由
	router.InitRouter(r)
	//grpc服务注册
	grpcServer := router.RegisterGrpc()
	stop := func() {
		grpcServer.Stop()
	}
	common.Run(r, config.AppConf.App.Name, config.AppConf.App.Addr, stop)
}
