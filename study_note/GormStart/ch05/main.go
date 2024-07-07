package main

import (
	"database/sql"
	"errors"
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

	//通过first查询单个数据, 获取第一条记录（主键升序）
	var user User
	//db.First(&user)

	//通过主键查询
	//db.First(&user, 10)
	// SELECT * FROM users WHERE id = 10;
	//我们不能给user赋值,然后执行主键查询
	// SELECT * FROM users WHERE id IN (1,2,3);
	result := db.First(&user, []int{1, 2, 3})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("未找到")
	}
	fmt.Println(user.ID)

	//检索全部对象
	var users []User
	result = db.Find(&users)
	fmt.Println("总共记录:", result.RowsAffected)
	for _, user := range users {
		fmt.Println(user.ID)
	}
}
