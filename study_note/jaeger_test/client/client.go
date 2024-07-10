package main

import (
	"context"
	"fmt"

	// 引入自定义的 OpenTracing gRPC 拦截器包
	"sale_master/study_note/jaeger_test/otgrpc"

	// 引入 OpenTracing 相关包
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	// 引入 gRPC 相关包
	"google.golang.org/grpc"

	// 引入 protobuf 生成的代码包
	"sale_master/study_note/jaeger_test/proto"
)

func main() {
	// Jaeger 配置
	cfg := jaegercfg.Configuration{
		// 指定采样器为常量采样器，即每个事务都会被采样
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		// 指定 Jaeger 的报告器配置，日志记录跨度
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.128.128:6831", // Jaeger Agent 的地址
		},
		// 服务名称
		ServiceName: "shop",
	}

	// 创建 Jaeger Tracer 实例
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(err)
	}
	// 设置全局 Tracer
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close() // 确保在程序退出时关闭 Tracer

	// gRPC 客户端连接
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())))
	if err != nil {
		panic(err)
	}
	defer conn.Close() // 确保在程序退出时关闭连接

	// 创建 gRPC 客户端
	c := proto.NewGreeterClient(conn)
	// 发起 gRPC 调用
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "llb"})
	if err != nil {
		panic(err)
	}
	// 输出调用结果
	fmt.Println(r.Message)
}
