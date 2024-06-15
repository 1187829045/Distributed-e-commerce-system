package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 返回一个 Gin 中间件函数，用于处理跨域请求
func Cors() gin.HandlerFunc {
	// 返回一个闭包函数作为 Gin 中间件处理函数
	return func(c *gin.Context) {
		method := c.Request.Method

		// 设置允许跨域的域名，"*" 表示允许所有域名访问，可以根据需求进行设置
		c.Header("Access-Control-Allow-Origin", "*")
		// 设置允许的请求头
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		// 设置允许的请求方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		// 设置浏览器可以访问的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		// 设置是否允许发送 Cookie（true 为允许）
		c.Header("Access-Control-Allow-Credentials", "true")

		// 如果请求方法为 OPTIONS，直接返回 StatusNoContent，表示允许预检请求通过
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 继续处理后续请求
		c.Next()
	}
}
