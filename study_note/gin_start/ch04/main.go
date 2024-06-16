package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义Person结构体
// 结构体字段包括ID和Name，并使用URI标签和绑定标签
type Person struct {
	ID   int    `uri:"id" binding:"required"`
	Name string `uri:"name" binding:"required"`
}

func main() {
	// 创建一个默认的Gin路由器，带有默认的中间件：日志和恢复中间件
	router := gin.Default()

	// 定义一个GET路由，带有动态URL参数:name和:id
	router.GET("/:name/:id", func(c *gin.Context) {
		// 创建一个Person结构体实例
		var person Person

		// 将URI中的参数绑定到Person结构体
		//这段代码的作用是从请求的URI中绑定参数到Person结构体实例，并处理绑定失败的情况
		if err := c.ShouldBindUri(&person); err != nil {
			// 如果绑定失败，返回404状态码
			c.Status(404)
			return
		}

		// 返回一个包含name和id的JSON对象
		c.JSON(http.StatusOK, gin.H{
			"name": person.Name,
			"id":   person.ID,
		})
	})

	// 在端口8083上启动Gin服务器
	router.Run(":8083")
}
