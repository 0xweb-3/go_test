package main

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go_test/jaeger_test/otgrpc"
	"go_test/jaeger_test/proto"
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

	// 设置 gRPC 服务器选项，使用 OpenTracing 拦截器
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())))

	g := grpc.NewServer(opts...)
	proto.RegisterHelloServer(g, &server{})
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err.Error())
	}
	g.Serve(lis)
}
