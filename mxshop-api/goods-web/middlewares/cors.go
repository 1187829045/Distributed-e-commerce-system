package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 返回一个 gin.HandlerFunc，它设置必要的 CORS（跨域资源共享）头部信息。
func Cors() gin.HandlerFunc {
	// 返回一个处理函数，用于处理每个请求。
	return func(c *gin.Context) {
		// 获取当前请求的方法（GET, POST, OPTIONS 等）
		method := c.Request.Method

		// 设置允许所有域名访问的头部信息
		c.Header("Access-Control-Allow-Origin", "*")
		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		// 设置可以被浏览器访问的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		// 设置是否允许浏览器发送 Cookie 和 HTTP 认证信息
		c.Header("Access-Control-Allow-Credentials", "true")

		// 如果是 OPTIONS 请求方法（通常是 CORS 预检请求），终止请求并返回 204 状态码
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
