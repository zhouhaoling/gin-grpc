package router

import (
	"github.com/gin-gonic/gin"
)

// Router 接口
type Router interface {
	Route(r *gin.Engine)
}

type RegisterRouter struct {
}

func NewRegisterRouter() *RegisterRouter {
	return &RegisterRouter{}
}

var routers []Router

func (rg *RegisterRouter) Route(ro Router, router *gin.Engine) {
	ro.Route(router)
}

func InitRouter(r *gin.Engine) {
	//rg := NewRegisterRouter()
	//rg.Route(user.NewRouterUser(), r)
	for _, ro := range routers {
		ro.Route(r)
	}
}

func Register(ro ...Router) {
	routers = append(routers, ro...)
}
