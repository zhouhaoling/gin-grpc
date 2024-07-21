package router

import (
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	mg "test.com/project-grpc/menu_grpc"

	auth "test.com/project-grpc/auth_grpc"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"test.com/common/discovery"
	"test.com/common/logs"
	ag "test.com/project-grpc/account_grpc"
	dg "test.com/project-grpc/department_grpc"
	pg "test.com/project-grpc/project_grpc"
	"test.com/project-project/config"
	"test.com/project-project/internal/rpc"
	"test.com/project-project/internal/service"
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
			pg.RegisterTaskServiceServer(server, service.NewTaskService())
			ag.RegisterAccountServiceServer(server, service.NewAccountService())
			dg.RegisterDepartmentServiceServer(server, service.NewDepartmentService())
			auth.RegisterAuthServiceServer(server, service.NewAuthService())
			mg.RegisterMenuServiceServer(server, service.NewMenuService())
		},
	}
	//cacheInterceptor := interceptor.NewCacheInterceptor()
	//server := grpc.NewServer(cacheInterceptor.Cache())
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			otelgrpc.UnaryServerInterceptor(),
			//interceptor.NewCacheInterceptor().CacheInterceptor(),
		)),
	)
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

func InitUserRpc() {
	rpc.InitGrpcUserClient()
}
