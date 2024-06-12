package project

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"test.com/common/discovery"
	"test.com/common/logs"
	"test.com/project-api/config"
	pg "test.com/project-grpc/project_grpc"
)

var projectServiceClient pg.ProjectServiceClient

func InitGrpcProjectClient() {
	etcdRegister := discovery.NewResolver(config.AppConf.Etcd.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	conn, err := grpc.Dial("etcd:///project", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	projectServiceClient = pg.NewProjectServiceClient(conn)
}
