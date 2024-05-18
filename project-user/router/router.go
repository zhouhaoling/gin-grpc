package router

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"test.com/project-user/config"
	"test.com/project-user/internal/dao"
	"test.com/project-user/internal/service"
	"test.com/project-user/internal/service/user_grpc"
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

type registerGrpc struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

func RegisterGrpc() *grpc.Server {
	c := registerGrpc{
		Addr: config.AppConf.Grpc.Addr,
		RegisterFunc: func(server *grpc.Server) {
			user_grpc.RegisterLoginServiceServer(server, service.NewUserService(dao.RC))
		},
	}
	server := grpc.NewServer()
	c.RegisterFunc(server)
	listen, err := net.Listen("tcp", config.AppConf.Grpc.Addr)
	if err != nil {
		log.Println("cannot listen")
	}

	go func() {
		err = server.Serve(listen)
		if err != nil {
			log.Println("server started error", err)
			return
		}
	}()

	return server
}
