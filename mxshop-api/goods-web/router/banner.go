package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/goods-web/api/banners"
	"mxshop-api/goods-web/middlewares"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners").Use(middlewares.Trace())
	{
		BannerRouter.GET("", banners.List)          // 轮播图列表页
		//依次使用 middlewares.JWTAuth() 和 middlewares.IsAdminAuth() 中间件进行
		//JWT 鉴权和管理员权限检查，然后执行 banners.Delete 处理函数，用于删除特定ID的轮播图。
		BannerRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Delete) // 删除轮播图
		BannerRouter.POST("",  middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.New)       //新建轮播图
		BannerRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Update) //修改轮播图信息
	}
}