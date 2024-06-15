package router

import (
	"github.com/gin-gonic/gin" // 导入 Gin 框架
	"mxshop-api/user-web/api"  // 导入用户相关的 API 处理函数
)

// InitBaseRouter 初始化基础路由
func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base") // 在传入的 RouterGroup 上创建一个 "base" 的路由组
	{
		// GET 请求：获取验证码，处理函数为 api.GetCaptcha
		BaseRouter.GET("captcha", api.GetCaptcha)

		// POST 请求：发送短信，处理函数为 api.SendSms
		BaseRouter.POST("send_sms", api.SendSms)
	}
}
