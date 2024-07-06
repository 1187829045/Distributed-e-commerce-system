package main

import (
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/uber/jaeger-client-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.128.128:6831",
		},
		ServiceName: "llb",
	}
	//生成链路
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	span := opentracing.StartSpan("go-grpc-web")
	time.Sleep(time.Second) //很重要
	defer span.Finish()
}
