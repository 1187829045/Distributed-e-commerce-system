package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/welcome", welcome)
	router.POST("/form_post", formPost)
	router.POST("/post", getPost)
	router.Run(":8083")
}

func getPost(c *gin.Context) {
	id := c.Query("id")
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

// 从post里面提取参数
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

// 从get里面获取参数
func welcome(c *gin.Context) {
	// 从查询参数中获取"firstname"，如果没有提供则默认值为"bobby"
	firstName := c.DefaultQuery("firstname", "llb")
	// 从查询参数中获取"lastname"，如果没有提供则默认值为"imooc"
	lastName := c.DefaultQuery("lastname", "cqupt")

	// 返回一个包含first_name和last_name的JSON对象
	c.JSON(http.StatusOK, gin.H{
		"first_name": firstName,
		"last_name":  lastName,
	})
}
