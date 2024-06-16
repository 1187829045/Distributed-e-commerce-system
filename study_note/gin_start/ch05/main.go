package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 创建一个默认的Gin路由器，带有默认的中间件：日志和恢复中间件
	router := gin.Default()

	// 定义一个GET路由，处理/welcome请求
	router.GET("/welcome", welcome)
	// 定义一个POST路由，处理/form_post请求
	router.POST("/form_post", formPost)
	// 定义另一个POST路由，处理/post请求
	router.POST("/post", getPost)

	// 在端口8083上启动Gin服务器
	router.Run(":8083")
}

// 处理/post请求的处理函数
func getPost(c *gin.Context) {
	// 从查询参数中获取"id"
	id := c.Query("id")
	// 从查询参数中获取"page"，如果没有提供则默认值为"0"
	page := c.DefaultQuery("page", "0")
	// 从表单参数中获取"name"
	name := c.PostForm("name")
	// 从表单参数中获取"message"，如果没有提供则默认值为"信息"
	message := c.DefaultPostForm("message", "信息")

	// 返回一个包含id, page, name, message的JSON对象
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"page":    page,
		"name":    name,
		"message": message,
	})
}

// 处理/form_post请求的处理函数
func formPost(c *gin.Context) {
	// 从表单参数中获取"message"
	message := c.PostForm("message")
	// 从表单参数中获取"nick"，如果没有提供则默认值为"anonymous"
	nick := c.DefaultPostForm("nick", "anonymous")

	// 返回一个包含message和nick的JSON对象
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"nick":    nick,
	})
}

// 处理/welcome请求的处理函数
func welcome(c *gin.Context) {
	// 从查询参数中获取"firstname"，如果没有提供则默认值为"bobby"
	firstName := c.DefaultQuery("firstname", "bobby")
	// 从查询参数中获取"lastname"，如果没有提供则默认值为"imooc"
	lastName := c.DefaultQuery("lastname", "imooc")

	// 返回一个包含first_name和last_name的JSON对象
	c.JSON(http.StatusOK, gin.H{
		"first_name": firstName,
		"last_name":  lastName,
	})
}
