package main

import (
	"context"
	_ "github.com/mbobakov/grpc-consul-resolver" // 这个必须引用
	"go_test/grpclb/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("consul://192.168.21.4:8500/user_srv?wait=14s&tag=user_srv",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	userSrvClient := proto.NewUserClient(conn)
	res, err := userSrvClient.GetUserList(context.Background(), &proto.GetUserListReq{
		PageSize:  10,
		PageToken: "1",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(res.GetUserInfos())
}
