package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func MyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// 设置一个示例的上下文变量，这里设置了一个名为 "example" 的键值对，值为 "123456"
		c.Set("example", "123456")
		// 让原本应该执行的后续逻辑继续执行
		c.Next()
		end := time.Since(t)
		fmt.Printf("耗时:%V\n", end)
		// 获取请求处理后的状态码
		status := c.Writer.Status()
		fmt.Println("状态", status)
	}
}

func Hook404() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
		status := c.Writer.Status()
		if status == 404 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "页面找不到",
			})
		}
	}
}

func main() {
	router := gin.Default()
	//使用logger和recovery中间件 全局所有
	router.Use(MyLogger(), Hook404())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run(":8083")
}
