package user

import (
	"github.com/gin-gonic/gin"
	"log"
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

func (ru *RouterUser) Route(r *gin.Engine) {
	h := NewHandlerUser()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
}
