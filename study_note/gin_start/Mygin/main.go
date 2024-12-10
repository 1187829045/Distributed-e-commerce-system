package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Phone    string `form:"phone" json:"phone" binding:"required"`
}

var user User

func MyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		time.Sleep(1 * time.Second)
		c.Next()
		end := time.Since(t)
		fmt.Printf("耗时:%d\n", end)
		// 获取请求处理后的状态码
		status := c.Writer.Status()
		fmt.Println("状态", status)
	}
}
func main() {

	r := gin.Default()

	r.Use(MyLogger())
	r.GET("/id/:id", func(c *gin.Context) {
		id := c.Param("id") //这是获取路劲参数
		c.JSON(200, gin.H{
			"message": "hello world",
			"id":      id,
		})
	})
	r.GET("/name", func(c *gin.Context) {
		name := c.Query("name") //这是获取查询参数
		c.JSON(200, gin.H{
			"message": "hello world",
			"name":    name,
		})
	})
	r.GET("/userinfo", func(c *gin.Context) {
		err := c.ShouldBind(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"message": err.Error(),
			})
		}
		c.JSON(200, gin.H{
			"name":     user.Username,
			"phone":    user.Phone,
			"password": user.Password,
		})
	})
	fmt.Println("servcer start port:8080 ")
	go func() {
		r.Run(":8080")
	}()

	// 创建一个通道，用于接收系统信号
	quit := make(chan os.Signal)

	// 监听系统信号：SIGINT（Ctrl+C 中断信号）和 SIGTERM（终止信号）
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号的到来，一旦接收到信号，程序将继续往下执行
	<-quit
	log.Println("Shutdown Server ...")
	time.Sleep(3 * time.Second)
	// 当接收到信号时，打印日志并执行优雅关闭服务器的逻辑
}
