package rpc

import (
	"log"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"test.com/common/discovery"
	"test.com/common/logs"
	"test.com/project-api/config"
	ug "test.com/project-grpc/user_grpc"
)

var LoginServiceClient ug.LoginServiceClient
var UsersServiceClient ug.UserServiceClient

func InitGrpcUserClient() {
	etcdRegister := discovery.NewResolver(config.AppConf.Etcd.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial("etcd:///user",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	LoginServiceClient = ug.NewLoginServiceClient(conn)
	UsersServiceClient = ug.NewUserServiceClient(conn)
}
