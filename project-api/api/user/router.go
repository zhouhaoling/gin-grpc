package user

import (
	"log"

	"test.com/project-api/api/middleware"
	"test.com/project-api/api/rpc"

	"github.com/gin-gonic/gin"
	"test.com/project-api/config"
	"test.com/project-api/router"
)

func init() {
	log.Println("init user router")
	router.Register(NewRouterUser())
}

type RouterUser struct {
}

func NewRouterUser() *RouterUser {
	return &RouterUser{}
}

// Route 用户路由模块
/*
 * @Description:存放用户相关路由
 */
func (ru *RouterUser) Route(r *gin.Engine) {
	rpc.InitGrpcUserClient()
	h := NewHandlerUser()
	pl := r.Group("/project/login")
	pl.Use(middleware.NewJwtMiddlewareBuilder().IgnorePaths(config.MemberPaths).Build())
	pl.POST("/getCaptcha", h.getCaptcha)
	pl.POST("", h.userLogin)
	pl.POST("/register", h.userRegister)
	org := r.Group("/project/organization")
	org.Use(middleware.NewJwtMiddlewareBuilder().Build())
	org.POST("/_getOrgList", h.myOrgList)
}
