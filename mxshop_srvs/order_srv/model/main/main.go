package main

import (
	"crypto/md5"                // 引入 md5 包，用于生成 MD5 哈希
	"encoding/hex"              // 引入 hex 包，用于将 MD5 哈希编码为十六进制字符串
	"gorm.io/driver/mysql"      // 引入 GORM 的 MySQL 驱动
	"gorm.io/gorm"              // 引入 GORM 包
	"gorm.io/gorm/logger"       // 引入 GORM 的 logger 包
	"gorm.io/gorm/schema"       // 引入 GORM 的 schema 包，用于定义数据库模式
	"io"                        // 引入 io 包，用于输入/输出操作
	"log"                       // 引入 log 包，用于日志记录
	"os"                        // 引入 os 包，用于操作系统功能
	"shop_srvs/order_srv/model" // 引入自定义的 model 包
	"time"                      // 引入 time 包，用于时间操作
)

// 生成 MD5 哈希值的函数
func genMd5(code string) string {
	Md5 := md5.New()                        // 创建一个新的 MD5 哈希对象
	_, _ = io.WriteString(Md5, code)        // 将字符串写入哈希对象
	return hex.EncodeToString(Md5.Sum(nil)) // 将哈希值转换为十六进制字符串并返回
}

func main() {
	// 数据库连接字符串
	dsn := "root:root@tcp(192.168.128.128:3306)/shop_order_srv?charset=utf8mb4&parseTime=True&loc=Local"

	// 配置新的日志记录器
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 设置日志输出到标准输出
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值为 1 秒
			LogLevel:      logger.Info, // 设置日志级别为 Info
			Colorful:      true,        // 启用彩色打印
		},
	)

	// 全局模式配置
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: newLogger, // 使用配置的日志记录器
	})
	if err != nil {
		panic(err) // 如果连接数据库失败，抛出错误
	}

	// 自动迁移数据库模式
	_ = db.AutoMigrate(&model.ShoppingCart{}, &model.OrderInfo{}, &model.OrderGoods{})
}
