package iniitialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-api/user-web/middlewares"
	router2 "shop-api/user-web/router"
)

// 函数初始化并返回一个带有预定义路由和中间件的 Gin 引擎。

func Routers() *gin.Engine {
	Router := gin.Default()

	// 定义一个健康检查端点，用于验证服务是否运行。
	Router.GET("/health", func(c *gin.Context) {
		// 返回 HTTP 200 状态码，并返回一个 JSON 响应表示成功。
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
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

	return Router
}
