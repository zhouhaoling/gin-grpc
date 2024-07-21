package main

import (
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"test.com/project-api/tracing"

	"github.com/gin-gonic/gin"
	srv "test.com/common"
	_ "test.com/project-api/api"
	"test.com/project-api/config"
	"test.com/project-api/router"
)

func main() {

	r := gin.Default()
	tp, tpErr := tracing.JaegerTraceProvider()
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	r.Use(otelgin.Middleware("project-api"))
	//静态目录映射
	r.StaticFS("/upload", http.Dir("upload"))
	//路由
	router.InitRouter(r)
	//开启pprof 默认访问路径是/debug/pprof, 可以自定义访问路径(r, "/dev/pprof")
	//pprof.Register(r)
	srv.Run(r, config.AppConf.App.Name, config.AppConf.App.Addr, nil)
}
