package user

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ug "test.com/project-user/user_grpc"
)

var LoginServiceClient ug.LoginServiceClient

func InitGrpcUserClient() {
	conn, err := grpc.NewClient("127.0.0.1:8881", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	LoginServiceClient = ug.NewLoginServiceClient(conn)
}
