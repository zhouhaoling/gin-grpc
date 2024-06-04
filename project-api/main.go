package main

import (
	"github.com/gin-gonic/gin"
	srv "test.com/common"
	_ "test.com/project-api/api/user"

	"test.com/project-api/config"
	"test.com/project-api/router"
)

func main() {

	r := gin.Default()
	//路由
	router.InitRouter(r)
	srv.Run(r, config.AppConf.App.Name, config.AppConf.App.Addr, nil)
}
