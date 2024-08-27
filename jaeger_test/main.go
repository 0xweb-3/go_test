package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"time"
)

func main() {
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
	opentracing.SetGlobalTracer(tracer) // 设置全局的
	defer closer.Close()

	parentSpan := tracer.StartSpan("main_span")

	span := tracer.StartSpan("go-grpc-web1", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 500)
	span.Finish()

	span1 := tracer.StartSpan("go-grpc-web2", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 1000)
	span1.Finish()

	span3 := opentracing.StartSpan("go-grpc-web3", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(time.Millisecond * 1500)
	span3.Finish()

	parentSpan.Finish()
}
