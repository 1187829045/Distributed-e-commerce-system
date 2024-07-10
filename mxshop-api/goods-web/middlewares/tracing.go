package middlewares

import (
	"fmt"
	"shop-api/goods-web/global"

	"github.com/gin-gonic/gin"              // 引入 Gin 框架
	"github.com/opentracing/opentracing-go" // 引入 OpenTracing 包
	"github.com/uber/jaeger-client-go"      // 引入 Jaeger 客户端包
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// Trace 是一个 Gin 中间件，用于实现分布式跟踪
func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 配置 Jaeger
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst, // 使用常量采样器
				Param: 1,                       // 采样率为 1，即每个事务都采样
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans: true, // 记录跨度日志
				// 设置 Jaeger Agent 的地址
				LocalAgentHostPort: fmt.Sprintf("%s:%d", global.ServerConfig.JaegerInfo.Host, global.ServerConfig.JaegerInfo.Port),
			},
			// 设置服务名称
			ServiceName: global.ServerConfig.JaegerInfo.Name,
		}

		// 创建 Jaeger Tracer 实例
		tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			panic(err)
		}
		// 设置全局 Tracer
		opentracing.SetGlobalTracer(tracer)
		defer closer.Close() // 确保在函数返回前关闭 Tracer

		// 在每个请求中创建一个新的 Span
		startSpan := tracer.StartSpan(ctx.Request.URL.Path)
		defer startSpan.Finish() // 确保在请求处理结束后完成 Span

		// 将 Tracer 和 Span 存储到 Gin 上下文中，以便后续中间件和处理函数使用
		ctx.Set("tracer", tracer)
		ctx.Set("parentSpan", startSpan)

		// 执行下一个中间件或处理函数
		ctx.Next()
	}
}
