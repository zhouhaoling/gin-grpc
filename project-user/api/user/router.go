package user

import (
	"log"

	"github.com/gin-gonic/gin"
	"test.com/project-user/router"
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
	//h := &HandlerUser{}
	h := NewHandlerUser()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
}
