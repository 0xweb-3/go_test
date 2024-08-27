package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go_test/jaeger_test/otgrpc"
	"go_test/jaeger_test/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func main() {
	// jaeger配置
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{ //采样
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true, // 是否打印日志
			LocalAgentHostPort: "192.168.21.3:6831",
		},
		ServiceName: "xinshop",
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer) // 将 tracer 设置为全局的
	defer closer.Close()

	// 客户端拦截器
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	// 这里添加内容
	//opts = append(opts, grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	opts = append(opts, grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
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
