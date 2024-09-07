package main

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Product struct {
	gorm.Model
	Code  sql.NullString
	Price uint
}

func main() {
	//charset=utf8mb4：指定字符集为 utf8mb4，这是 UTF-8 的一个超集，支持存储 Emoji 字符等 4 字节字符。确保数据在存储和读取时不会出现字符编码问题。
	//parseTime=True：将 MySQL 中的 DATETIME、DATE、TIMESTAMP 字段自动解析为 Go 语言的 time.Time 类型。如果不设置，时间类型可能会被解析为字符串。
	//loc=Local：设置时区为本地时间。MySQL 中的时间通常以 UTC 存储，使用 loc=Local 可以让 Go 语言中的 time.Time 类型自动转换为本地时区。
	dsn := "数据库用户名:数据库用户密码@tcp(192.168.128.128:3306)/数据库的名称?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	//定义一个表结构， 将表结构直接生成对应的表 - migrations
	_ = db.AutoMigrate(&Product{})
	db.Create(&Product{Code: sql.NullString{"D42", true}, Price: 100})
	var product Product
	db.First(&product, 1)                 // 根据整形主键查找
	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 100)
	// Update - 更新多个字段
	db.Model(&product).Updates(Product{Price: 200, Code: sql.NullString{"", true}}) // 仅更新非零值字段,所以引入sql.NullString
	//如果我们去更新一个product 只设置了price：200
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	// Delete - 删除 product， 并没有执行delete语句，逻辑删除
	db.Delete(&product, 1)
}
