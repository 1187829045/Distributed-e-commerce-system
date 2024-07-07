package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivedAt    sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// 批量插入
func main() {
	dsn := "root:root@tcp(192.168.128.128:3306)/gorm_test?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	//单一的 SQL 语句
	var users = []User{{Name: "llb1"}, {Name: "llb2"}, {Name: "llb3"}}
	//db.Create(&users)

	//为什么不一次性提交所有的 还要分批次，sql语句有长度限制
	db.CreateInBatches(users, 10000) //第二个参数表示每次提交的最多数量

	for _, user := range users {
		fmt.Println(user.ID) // 1,2,3
	}

	db.Model(&User{}).Create(map[string]interface{}{
		"Name": "llb", "Age": 18,
	})
}
