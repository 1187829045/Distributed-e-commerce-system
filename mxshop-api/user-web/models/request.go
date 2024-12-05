package models

import (
	"github.com/dgrijalva/jwt-go"
)

// 定义了自定义的 JWT 载荷结构体

type CustomClaims struct {
	// 用户ID
	ID uint `json:"id"`
	// 用户昵称
	NickName string `json:"nickname"`
	// 权限ID，用于标识用户的权限级别
	AuthorityId uint `json:"authority_id"`
	// jwt.StandardClaims 是 JWT-go 库中预定义的标准声明结构体，包含了 JWT 的标准字段（如过期时间等）
	jwt.StandardClaims
}
