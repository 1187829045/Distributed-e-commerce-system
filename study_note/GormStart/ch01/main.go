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
	dsn := "root:root@tcp(192.168.128.128:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"
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
