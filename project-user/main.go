package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"test.com/common"
	"test.com/common/logs"
	"test.com/project-user/config"
	"test.com/project-user/internal/repository/database"
	"test.com/project-user/pkg/snowflake"
	"test.com/project-user/router"
)

func init() {
	//if err := config.InitConfig(); err != nil {
	//	log.Fatalln("初始化配置文件失败", err)
	//}
	if err := database.InitRedis(config.AppConf); err != nil {
		log.Fatalln("初始化redis失败", err)
	}
	if err := database.InitMySQL(); err != nil {
		log.Fatalln("初始化mysql失败", err)
	}
	if err := snowflake.InitSnowflake(config.AppConf.Snow.StartTime, config.AppConf.Snow.MachineID); err != nil {
		log.Fatalln("初始化雪花算法失败", err)
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
