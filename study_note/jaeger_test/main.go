package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"time"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	// 配置 Jaeger 客户端
	cfg := jaegercfg.Configuration{
		// 采样配置
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, // 采样类型，常量采样
			Param: 1,                       // 采样参数，这里表示全部采样
		},
		// 报告配置
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,                   // 是否记录 span 日志
			LocalAgentHostPort: "192.168.128.128:6831", // Jaeger agent 的地址
		},
		ServiceName: "llb", // 服务名称
	}

	// 生成 tracer 和 closer
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err) // 如果生成 tracer 失败，则抛出异常
	}

	// 设置全局 tracer
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close() // 在 main 函数结束前关闭 tracer

	// 开始一个新的 span
	span := opentracing.StartSpan("go-grpc-web")
	time.Sleep(time.Second) // 模拟操作延迟，确保 span 能够被记录
	defer span.Finish()     // 在 span 结束前关闭它
}
