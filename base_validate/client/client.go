package main

import (
	"context"
	"fmt"
	"go_test/base_validate/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func main() {
	// 客户端拦截器
	interceptor := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		fmt.Println(time.Since(start))
		return err
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(interceptor))

	conn, err := grpc.Dial("127.0.0.1:8080", opts...)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	c := proto.NewHelloClient(conn)
	r, err := c.Hello(context.Background(), &proto.HelloRequest{ // 这里使用自己定义的ctx
		Name:    "xin",
		AddTime: timestamppb.New(time.Now()),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Message)
}
