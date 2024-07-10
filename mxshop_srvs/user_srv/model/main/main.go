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
	////将用户密码变一下 随机字符串+用户密码
	//fmt.Println(genMd5("xxxxx_123456"))
	////使用默认选项生成带盐的编码密码
	//salt, encodedPwd := password.Encode("generic password", nil)
	//fmt.Println(salt, encodedPwd)
	//// 验证输入的密码是否与编码后的密码匹配
	//check := password.Verify("generic password", salt, encodedPwd, nil)
	//fmt.Println(check) // 打印验证结果，输出 true 表示验证通过

	// 使用自定义选项生成带盐的编码密码
	options := &password.Options{10, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("123456", options)
	newpassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newpassword)

	// 验证输入的密码是否与使用相同选项编码后的密码匹配
	passwordInfo := strings.Split(newpassword, "$")
	fmt.Println(passwordInfo)
	check := password.Verify("123456", passwordInfo[2], passwordInfo[3], options)
	fmt.Println(check) // 打印验证结果，输出 true 表示验证通过
}
