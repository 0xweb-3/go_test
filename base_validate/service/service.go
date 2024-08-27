package main

import (
	"context"
	"fmt"
	"go_test/base_validate/proto"
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
	//  拦截器代码
	interceptor := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		fmt.Println("收到一个新的请求")
		return handler(ctx, req)
	}
	opt := grpc.UnaryInterceptor(interceptor)
	g := grpc.NewServer(opt) // 拦截器注入
	proto.RegisterHelloServer(g, &server{})
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err.Error())
	}
	g.Serve(lis)
}
