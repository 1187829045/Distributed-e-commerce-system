package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"io"
	"strings"
)

// 生成MD5哈希值的函数
func genMd5(code string) string {
	// 创建一个新的MD5哈希对象
	Md5 := md5.New()
	// 将输入字符串写入到MD5哈希对象中进行处理
	// io.WriteString 返回两个值（写入的字节数和错误信息），这里用下划线 "_" 忽略它们
	_, _ = io.WriteString(Md5, code)
	// 计算输入字符串的MD5哈希值，并将其转换为十六进制字符串后返回
	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	/*
		dsn := "root:root@tcp(192.168.128.134:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

			newLogger := logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold: time.Second, // 慢 SQL 阈值
					LogLevel:      logger.Info, // Log level
					Colorful:      true,        // 禁用彩色打印
				},
			)

			// 全局模式
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
				Logger: newLogger,
			})
			if err != nil {
				panic(err)
			}
			_ = db.AutoMigrate(&model.User{})
	*/
	//fmt.Println(genMd5("xxxxx_123456"))
	//将用户密码变一下 随机字符串+用户密码
	//e10adc3949ba59abbe56e057f20f883e

	//使用默认选项生成带盐的编码密码
	salt, encodedPwd := password.Encode("generic password", nil)
	//fmt.Println(salt, encodedPwd)
	//// 验证输入的密码是否与编码后的密码匹配
	//check := password.Verify("generic password", salt, encodedPwd, nil)
	//fmt.Println(check) // 打印验证结果，输出 true 表示验证通过
	//

	// 使用自定义选项生成带盐的编码密码
	options := &password.Options{10, 100, 32, sha512.New}
	salt, encodedPwd = password.Encode("generic password", options)
	newpassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newpassword)
	// 验证输入的密码是否与使用相同选项编码后的密码匹配
	passwordInfo := strings.Split(newpassword, "$")
	fmt.Println(passwordInfo)
	check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
	fmt.Println(check) // 打印验证结果，输出 true 表示验证通过
}
