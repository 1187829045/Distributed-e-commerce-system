package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 表单验证
// 定义一个全局的翻译器变量
var trans ut.Translator

// 定义一个用于登录的表单结构体

type LoginForm struct {
	User     string `json:"user" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required"`
}

// 定义一个用于注册的表单结构体

type SignUpForm struct {
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
	Name       string `json:"name" binding:"required,min=3"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"` // 跨字段验证，确保与Password字段相同
}

// 去除结构体字段名中的前缀
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// 初始化翻译器

func InitTrans(locale string) (err error) {
	// 修改gin框架中的validator引擎属性，实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		//RegisterTagNameFunc 函数是在初始化验证器时调用的，它用于注册一个自定义的方法，该方法用于从结构体字段的标签中提取字段名。
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		// 第一个参数是备用的语言环境，后面的参数是支持的语言环境
		uni := ut.New(enT, zhT, enT)
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}

		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(v, trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, trans)
		default:
			en_translations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

func main() {
	// 初始化翻译器（设置为中文）
	if err := InitTrans("zh"); err != nil {
		fmt.Println("初始化翻译器错误")
		return
	}

	// 创建一个默认的Gin路由器，带有默认的中间件：日志和恢复中间件
	router := gin.Default()

	// 定义一个POST路由，处理/loginJSON请求
	router.POST("/loginJSON", func(c *gin.Context) {
		var loginForm LoginForm
		// 绑定请求数据到loginForm结构体，并进行验证
		if err := c.ShouldBind(&loginForm); err != nil {
			// 如果绑定或验证失败，处理错误
			//err 是一个错误（error）类型的变量，可能是任何实现了 error 接口的类型。
			//validator.ValidationErrors 是来自 go-playground/validator 包的一个类型，表示验证错误的集合。
			//通过类型断言，试图将 err 转换为 validator.ValidationErrors 类型
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				c.JSON(http.StatusOK, gin.H{
					"msg": err.Error(),
				})
				return
			}
			// 返回验证错误信息
			c.JSON(http.StatusBadRequest, gin.H{
				"error": removeTopStruct(errs.Translate(trans)),
			})
			return
		}
		// 返回登录成功信息
		c.JSON(http.StatusOK, gin.H{
			"msg": "登录成功",
		})
	})

	// 定义一个POST路由，处理/signup请求
	router.POST("/signup", func(c *gin.Context) {
		var signUpFrom SignUpForm
		// 绑定请求数据到signUpForm结构体，并进行验证
		if err := c.ShouldBind(&signUpFrom); err != nil {
			// 如果绑定或验证失败，返回错误信息
			fmt.Println(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// 返回注册成功信息
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册成功",
		})
	})

	// 在端口8083上启动Gin服务器
	_ = router.Run(":8083")
}
