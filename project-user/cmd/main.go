package main

import (
	"github.com/gin-gonic/gin"
	"test.com/common"
	_ "test.com/project-user/api/user"
	"test.com/project-user/router"
)

func main() {
	r := gin.Default()
	//路由
	router.InitRouter(r)
	common.Run(r, "project-user", ":8080")
}
