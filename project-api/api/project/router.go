package project

import (
	"log"

	"github.com/gin-gonic/gin"
	"test.com/project-api/router"
)

func init() {
	log.Println("init project router")
	router.Register(NewRouterProject())
}

type RouterProject struct {
}

func NewRouterProject() *RouterProject {
	return &RouterProject{}
}

// Route 项目路由模块
/*
 * @Description:存放用户相关路由
 */
func (ru *RouterProject) Route(r *gin.Engine) {
	InitGrpcProjectClient()
	p := NewHandlerProject()
	r.POST("/project/index", p.index)
}
