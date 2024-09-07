package global

import (
	"gorm.io/gorm"
	"shop_srvs/user_srv/config"
)

// 定义全局变量
var (
	DB           *gorm.DB            // GORM 的数据库连接对象
	ServerConfig config.ServerConfig // 服务器配置
	NacosConfig  config.NacosConfig  // Nacos 配置
)

// init 函数会在包初始化时自动执行
// 初始化数据库
//func init() {
//	// 定义数据库连接字符串（DSN），包含用户名、密码、主机地址、数据库名称和一些参数
//	dsn := "root:root@tcp(192.168.128.128:3306)/shop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
//
//	// 创建新的 GORM 日志对象
//	newLogger := logger.New(
//		// 设置日志输出到标准输出，带有时间戳的日志格式
//		log.New(os.Stdout, "\r\n", log.LstdFlags),
//		logger.Config{
//			SlowThreshold: time.Second, // 设置慢 SQL 阈值为 1 秒
//			LogLevel:      logger.Info, // 设置日志级别为 Info
//			Colorful:      true,        // 启用彩色打印
//		},
//	)
//
//	// 全局模式，使用 GORM 打开数据库连接
//	var err error
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true, // 禁用表名复数
//		},
//		Logger: newLogger, // 使用自定义的日志配置
//	})
//	// 如果连接数据库时发生错误，抛出异常
//	if err != nil {
//		panic(err)
//	}
//}
