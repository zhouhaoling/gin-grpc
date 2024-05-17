package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"test.com/common"
	"test.com/common/logs"
	_ "test.com/project-user/api/user"
	"test.com/project-user/router"
)

func main() {
	r := gin.New()
	r.Use(logs.GinLogger(), logs.GinRecovery(true))
	lc := &logs.LogConfig{
		DebugFileName: "gin-grpc\\common\\logs\\debug\\debug.log",
		InfoFileName:  "gin-grpc\\common\\logs\\info\\info.log",
		WarnFileName:  "gin-grpc\\common\\logs\\warn\\warn.log",
		MaxSize:       500,
		MaxAge:        28,
		MaxBackups:    3,
	}
	if err := logs.InitLogger(lc); err != nil {
		log.Fatalln("日志初始化失败")
	}

	//路由
	router.InitRouter(r)
	common.Run(r, "project-user", ":8080")
}
