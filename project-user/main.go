package main

import (
	"log"

	"test.com/project-user/config"
	"test.com/project-user/pkg/snowflake"

	"github.com/gin-gonic/gin"
	"test.com/common"
	"test.com/common/logs"
	_ "test.com/project-user/internal/repository/database"
	"test.com/project-user/router"
)

func main() {

	if err := snowflake.InitSnowflake(config.AppConf.Snow.StartTime, config.AppConf.Snow.MachineID); err != nil {
		log.Fatalln("初始化雪花算法失败", err)
	}
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
