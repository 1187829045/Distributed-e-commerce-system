package iniitialize

import (
	"github.com/gin-gonic/gin"           // 导入 Gin 框架用于构建 HTTP Web 应用程序
	"mxshop-api/user-web/middlewares"    // 导入 middlewares 包用于额外的请求处理
	router2 "mxshop-api/user-web/router" // 导入 router 包用于定义路由
	"net/http"                           // 导入 net/http 包用于 HTTP 状态码和服务器实现
)

// Routers 函数初始化并返回一个带有预定义路由和中间件的 Gin 引擎。
func Routers() *gin.Engine {
	Router := gin.Default() // 创建一个默认的 Gin 引擎实例，包含默认的中间件（日志记录和恢复）

	// 定义一个健康检查端点，用于验证服务是否运行。
	Router.GET("/health", func(c *gin.Context) {
		// 返回 HTTP 200 状态码，并返回一个 JSON 响应表示成功。
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK, // HTTP 200 状态码
			"success": true,          // 成功标志设置为 true
		})
	})

	// 使用 CORS 中间件处理跨源资源共享。
	Router.Use(middlewares.Cors())

	// 创建一个 API 分组，基础路径为 "/u/v1"。
	ApiGroup := Router.Group("/u/v1")

	// 在 API 分组中初始化与用户相关的路由。
	router2.InitUserRouter(ApiGroup)

	// 在 API 分组中初始化基础路由。
	router2.InitBaseRouter(ApiGroup)

	// 返回配置好的 Gin 引擎实例。
	return Router
}
