package main

import (
	"context"
	"fmt"
	"go_test/base/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func main() {
	// 客户端拦截器
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("127.0.0.1:8080", opts...)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	c := proto.NewHelloClient(conn)
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 1*time.Second)
	r, err := c.Hello(ctx, &proto.HelloRequest{ // 这里使用自己定义的ctx
		Name:    "xin",
		AddTime: timestamppb.New(time.Now()),
	})
	if err != nil {
		// 获取服务端抛出错误
		st, ok := status.FromError(err)
		if !ok {
			// 非grpc库中定义的错误
			panic(err)
		}
		fmt.Println("错误信息", st.Code(), "-===-", st.Message())
	}
	fmt.Println(r.Message)
}
