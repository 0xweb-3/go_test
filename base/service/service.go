package main

import (
	"context"
	"go_test/base/proto"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	proto.UnimplementedHelloServer // 只是起兼容作用，写法固定
}

func (s *server) Hello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: "利好xin====" + request.GetName(),
	}, nil
}

func main() {
	g := grpc.NewServer()
	proto.RegisterHelloServer(g, &server{})
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err.Error())
	}
	g.Serve(lis)
}
