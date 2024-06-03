package user

import (
	"log"

	"google.golang.org/grpc/resolver"

	"test.com/common/discovery"
	"test.com/common/logs"
	"test.com/project-api/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ug "test.com/project-user/user_grpc"
)

var LoginServiceClient ug.LoginServiceClient

func InitGrpcUserClient() {
	etcdRegister := discovery.NewResolver(config.AppConf.Etcd.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial("etcd:///user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	LoginServiceClient = ug.NewLoginServiceClient(conn)
}
