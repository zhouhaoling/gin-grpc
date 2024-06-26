package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"test.com/common"
	"test.com/common/logs"
	"test.com/project-project/config"
	"test.com/project-project/internal/repository/database"
	"test.com/project-project/router"
)

func init() {
	if err := database.InitRedis(config.AppConf); err != nil {
		log.Fatalln("初始化redis失败", err)
	}
}

func main() {

	r := gin.New()
	r.Use(logs.GinLogger(), logs.GinRecovery(true))
	//config.AppConf.InitZapLog()
	//zap.L().Info("初始化成功")
	//路由
	router.InitRouter(r)
	//grpc服务注册
	grpcServer := router.RegisterGrpc()
	//grpc服务注册到etcd
	router.RegisterEtcdServer()
	stop := func() {
		grpcServer.Stop()
	}
	common.Run(r, config.AppConf.App.Name, config.AppConf.App.Addr, stop)
}
