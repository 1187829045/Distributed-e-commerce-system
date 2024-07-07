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
	router := gin.Default()

	router.GET("/:name/:id", func(c *gin.Context) {
		var person Person
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
