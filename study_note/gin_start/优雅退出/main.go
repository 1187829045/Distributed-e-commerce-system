package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 创建一个默认的 Gin 路由器实例
	router := gin.Default()

	// 定义一个处理 GET 请求的路由，当访问 "/" 路径时，返回 JSON 响应
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "pong", // 返回一个包含 "msg": "pong" 的 JSON 响应，HTTP 状态码为 200
		})
	})

	// 启动一个新的 goroutine 来运行 Gin 服务器，监听在 8083 端口
	go func() {
		router.Run(":8083")
	}()

	// 创建一个通道，用于接收系统信号
	quit := make(chan os.Signal)

	// 监听系统信号：SIGINT（Ctrl+C 中断信号）和 SIGTERM（终止信号）
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待信号的到来，一旦接收到信号，程序将继续往下执行
	<-quit

	// 当接收到信号时，打印日志并执行优雅关闭服务器的逻辑
	log.Println("Shutdown Server ...")
}
