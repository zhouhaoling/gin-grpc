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

// routers里面存放这各个路由模块的对象，比如user.RouterUser
var routers []Router

func (rg *RegisterRouter) Route(ro Router, router *gin.Engine) {
	ro.Route(router)
}

func InitRouter(r *gin.Engine) {
	//rg := NewRegisterRouter()
	//rg.Route(&user.RouterUser{}, r)
	//遍历各个路由模块对象，并调用它们的路由注册方法
	for _, ro := range routers {
		ro.Route(r)
	}
}

func Register(ro ...Router) {
	routers = append(routers, ro...)
}
