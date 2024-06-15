package main

import (
	"fmt" // 导入 fmt 包用于格式化 I/O
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap" // 导入 zap 日志库
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/iniitialize" // 导入自定义包，用于初始化路由
	myvalidator "mxshop-api/user-web/validator"
)

func main() {
	// 1. 初始化日志系统
	iniitialize.InitLogger()

	// 2. 初始化配置文件
	iniitialize.InitConfig()

	// 3. 初始化路由
	Router := iniitialize.Routers()

	// 4. 初始化翻译功能，指定使用中文翻译
	if err := iniitialize.InitTrans("zh"); err != nil {
		// 如果翻译初始化失败，打印错误信息并终止程序
		fmt.Println("翻译失败")
		panic(err)
	}

	// 注册验证器
	//用于获取当前的验证引擎实例
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// v 是获取到的验证引擎实例，这里将自定义的手机号验证函数 ValidateMobile 注册到引擎中
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)

		// 注册手机号验证的错误信息翻译
		_ = v.RegisterTranslation("mobile", global.Trans,
			// 添加手机号验证错误的翻译信息到翻译器中，"{0} 非法的手机号码!" 是错误消息的模板，{0} 会被实际的字段值替代
			func(ut ut.Translator) error {
				return ut.Add("mobile", "{0} 非法的手机号码!", true)
			},
			// 定义翻译器如何生成具体的错误消息，fe.Field() 获取错误字段的名称
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("mobile", fe.Field())
				return t
			})
	}

	// zap.S() 是 zap 日志库中的一个全局函数，用于返回全局的 SugaredLogger 实例
	// zap.S().Debugf("启动服务器，端口:%d", port) 这行代码使用 zap 日志库记录了一条调试级别的日志信息。
	zap.S().Debugf("启动服务器，端口:%d", global.ServerConfig.Port)

	// 启动 Gin 服务器，监听并运行在指定端口上
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		// 如果服务器启动失败，使用 zap 记录错误日志并终止程序
		zap.S().Panic("启动失败", zap.Error(err))
	}
}
