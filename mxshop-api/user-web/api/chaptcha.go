package api

import (
	"github.com/gin-gonic/gin"        // 引入Gin框架，用于处理HTTP请求和响应
	"github.com/mojocn/base64Captcha" // 引入base64Captcha库，用于生成验证码
	"go.uber.org/zap"                 // 引入zap库，用于日志记录
	"net/http"                        // 引入net/http包，用于HTTP状态码
)

// var store = base64Captcha.DefaultMemStore 初始化一个默认的内存存储，用于存储验证码
var store = base64Captcha.DefaultMemStore

// gin.Context 是Gin框架中非常重要的一个结构体，它封装了处理 HTTP 请求的上下文信息，包括请求和响应的数据、参数、查询字符串、表单数据等。
// 通过 gin.Context，开发者可以方便地获取和操作 HTTP 请求和响应的数据，处理中间件逻辑，进行错误处理，控制请求的流程等等。

// GetCaptcha 处理获取验证码的请求
func GetCaptcha(ctx *gin.Context) {
	// 创建一个数字验证码的驱动，设置验证码的高度、宽度、位数、干扰强度、最大值
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	// 使用上面的驱动和内存存储创建一个新的验证码对象
	cp := base64Captcha.NewCaptcha(driver, store)
	// 生成验证码，返回验证码ID和base64编码的图像字符串
	id, b64s, _, err := cp.Generate()
	// 检查是否有错误
	if err != nil {
		zap.S().Errorf("生成验证码错误: ", err.Error()) // 如果有错误，记录错误日志
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成验证码错误",
		})
		return // 终止函数执行
	}
	// 如果没有错误，返回HTTP 200状态码和验证码ID、图像路径
	ctx.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}
