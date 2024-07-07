package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	goodsGroup := router.Group("/goods")
	{
		goodsGroup.GET("", goodsList)

		goodsGroup.GET("/:id/:action/add", goodsDetail)

		goodsGroup.POST("", createGoods)
	}
	// 在端口8083上启动Gin服务器
	router.Run(":8083")
}

// 创建新商品的处理函数
// 这是一个占位函数；你可以在这里添加创建商品的逻辑
func createGoods(c *gin.Context) {
	// 目前，这个函数什么也不做
}

// 获取特定商品详细信息的处理函数
func goodsDetail(c *gin.Context) {
	// 从URL中提取'id'参数
	id := c.Param("id")

	// 从URL中提取'action'参数
	action := c.Param("action")

	// 返回一个包含id和action的JSON对象
	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"action": action,
	})
}

// 列出所有商品的处理函数
func goodsList(context *gin.Context) {
	// 返回一个包含简单消息的JSON对象
	context.JSON(http.StatusOK, gin.H{
		"name": "goodsList",
	})
}
