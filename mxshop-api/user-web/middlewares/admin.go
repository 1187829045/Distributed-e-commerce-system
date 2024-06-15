package middlewares

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/models"
	"net/http"
)

// IsAdminAuth 返回一个 Gin 中间件函数，用于判断当前用户是否为管理员用户。
func IsAdminAuth() gin.HandlerFunc {
	// 返回一个闭包函数作为 Gin 中间件处理函数
	return func(c *gin.Context) {
		// 从 Gin 上下文中获取声明的认证信息 claims
		claims, _ := c.Get("claims")
		// 将 claims 转换为自定义声明结构体 *models.CustomClaims
		currentUser := claims.(*models.CustomClaims)
		// 判断当前用户的权限 ID 是否不等于 2（假设 2 是管理员的权限 ID）
		if currentUser.AuthorityId != 2 {
			// 如果不是管理员，返回 HTTP 状态码 403 Forbidden 和 JSON 格式的错误消息
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			// 终止请求处理流程，后续的中间件和处理函数将不再执行
			c.Abort()
			return
		}
		// 如果是管理员，调用 c.Next() 继续执行后续的中间件和处理函数
		c.Next()
	}
}
