package project

import (
	"log"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"test.com/common/discovery"
	"test.com/common/logs"
	"test.com/project-api/config"
	ag "test.com/project-grpc/account_grpc"
	auth "test.com/project-grpc/auth_grpc"
	dg "test.com/project-grpc/department_grpc"
	mg "test.com/project-grpc/menu_grpc"
	pg "test.com/project-grpc/project_grpc"
)

var projectServiceClient pg.ProjectServiceClient
var taskServiceClient pg.TaskServiceClient
var AccountServiceClient ag.AccountServiceClient
var departmentServiceClient dg.DepartmentServiceClient
var authServiceClient auth.AuthServiceClient
var menuServiceClient mg.MenuServiceClient

func InitGrpcProjectClient() {
	etcdRegister := discovery.NewResolver(config.AppConf.Etcd.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial("etcd:///project",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	projectServiceClient = pg.NewProjectServiceClient(conn)
	taskServiceClient = pg.NewTaskServiceClient(conn)
	AccountServiceClient = ag.NewAccountServiceClient(conn)
	departmentServiceClient = dg.NewDepartmentServiceClient(conn)
	authServiceClient = auth.NewAuthServiceClient(conn)
	menuServiceClient = mg.NewMenuServiceClient(conn)
}
