package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	srv "test.com/common"
	"test.com/common/logs"
	"test.com/project-project/config"
	"test.com/project-project/router"
	"test.com/project-project/tracing"
)

//func init() {
//	if err := database.InitRedis(config.AppConf); err != nil {
//		log.Fatalln("初始化redis失败", err)
//	}
//	if err := database.InitMySQL(); err != nil {
//		log.Fatalln("初始化mysql失败", err)
//	}
//}

func main() {

	r := gin.New()
	tp, tpErr := tracing.JaegerTraceProvider()
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	r.Use(logs.GinLogger(), logs.GinRecovery(true))
	//config.AppConf.InitZapLog()
	//zap.L().Info("初始化成功")
	//路由
	router.InitRouter(r)
	router.InitUserRpc()
	//grpc服务注册
	grpcServer := router.RegisterGrpc()
	//grpc服务注册到etcd
	router.RegisterEtcdServer()
	stop := func() {
		grpcServer.Stop()
	}
	srv.Run(r, config.AppConf.App.Name, config.AppConf.App.Addr, stop)
}
