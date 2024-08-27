package main

import (
	"context"
	"fmt"
	"go_test/grpc_err/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

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
