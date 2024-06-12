package router

import (
	"log"
	"net"

	pg "test.com/project-grpc/project_grpc"
	"test.com/project-project/internal/service"

	"google.golang.org/grpc/resolver"

	"test.com/common/logs"

	"test.com/common/discovery"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"test.com/project-project/config"
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

// RegisterGrpc 注册grpc服务
func RegisterGrpc() *grpc.Server {
	c := registerGrpc{
		Addr: config.AppConf.Grpc.Addr,
		RegisterFunc: func(server *grpc.Server) {
			pg.RegisterProjectServiceServer(server, service.NewProjectService())
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

// RegisterEtcdServer 注册etcd服务
func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(config.AppConf.Etcd.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	info := discovery.Server{
		Name:    config.AppConf.Grpc.Name,
		Addr:    config.AppConf.Grpc.Addr,
		Version: config.AppConf.Grpc.Version,
		Weight:  config.AppConf.Grpc.Weight,
	}
	r := discovery.NewRegister(config.AppConf.Etcd.Addrs, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}
