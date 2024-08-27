package main

import (
	"context"
	"fmt"
	"go_test/grpc_auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
	interceptor := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		var userId string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			userId = md.Get("userId")[0]
			fmt.Println("用户id", userId)
		} else {
			return resp, status.Error(codes.Unauthenticated, "获取用户id失败") // 使用grpc错误
		}

		if userId != "00001" {
			return resp, status.Error(codes.Unauthenticated, "认证失败")
		}

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
