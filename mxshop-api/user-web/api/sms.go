package api

import (
	"context"
	"fmt"
	"math/rand"
	"mxshop-api/user-web/forms"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"      // 导入阿里云 SDK 的请求包
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi" // 导入阿里云短信服务包
	"github.com/gin-gonic/gin"                                 // 导入 Gin 框架
	"github.com/go-redis/redis/v8"                             // 导入 Redis 客户端
	"mxshop-api/user-web/global"                               // 导入全局配置包
)

// GenerateSmsCode 生成指定长度的随机数字短信验证码
func GenerateSmsCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} // 数字数组
	r := len(numeric)
	// 使用当前时间的 UnixNano 作为随机数种子
	//UnixNano 是 Go 语言中用于获取当前时间的函数之一，它返回的是当前时间的纳秒级别的 Unix 时间戳。
	rand.Seed(time.Now().UnixNano())
	//strings.Builder 是 Go 语言中提供的一个用于高效构建字符串的类型。
	//它在 Go 1.10 版本中引入，用于替代传统的字符串拼接方式（如使用 + 或 fmt.Sprintf）来避免因字符串拼接导致的性能问题。
	var sb strings.Builder
	for i := 0; i < width; i++ {
		// 随机生成数字并添加到字符串构建器中
		//因此，rand.Intn(r) 返回的随机数可以是 0 到 r-1 之间的任意整数。
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String() // 返回生成的随机数字字符串
}

// SendSms 处理发送短信的请求
func SendSms(ctx *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	// 绑定表单数据到结构体，并检查是否有错误
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		// 处理验证器错误
		HandleValidatorError(ctx, err)
		return
	}

	// 使用阿里云 SDK 创建短信服务客户端
	client, err :=
		dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey,
			global.ServerConfig.AliSmsInfo.ApiSecrect)

	if err != nil {
		// 出现错误时抛出异常
		panic(err)
	}

	// 生成 6 位数字的短信验证码
	smsCode := GenerateSmsCode(6)

	// 创建一个通用请求对象
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // 使用 HTTPS 协议
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile            // 手机号
	request.QueryParams["SignName"] = "慕学在线"                        // 阿里云验证过的项目名，自行设置
	request.QueryParams["TemplateCode"] = "SMS_181850725"               // 阿里云的短信模板号，自行设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" // 短信模板中的验证码内容，自动生成

	// 发送短信请求并获取响应
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response)) // 打印请求响应信息
	if err != nil {
		fmt.Print(err.Error()) // 打印错误信息
	}

	// 将验证码保存到 Redis 中，使用手机号作为键名，并设置过期时间
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port), // Redis 服务器地址和端口
	})
	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)

	// 返回发送成功的 JSON 响应
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
