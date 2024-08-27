package main

import (
	"context"
	"go_test/grpc_err/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type server struct {
	proto.UnimplementedHelloServer // 只是起兼容作用，写法固定
}

func (s *server) Hello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloResponse, error) {
	return nil, status.Error(codes.InvalidArgument, "invalid userId")
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
