package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"shop-api/user-web/global"
	"shop-api/user-web/models"
	"time"
)

// 是一个 gin 中间件函数，用于 JWT 认证

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头部获取 token 信息，头部字段名称为 x-token
		token := c.Request.Header.Get("x-token")
		if token == "" {
			// 如果没有 token，返回 401 未授权状态码和提示信息
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort() // 中止请求
			return
		}

		j := NewJWT() // 创建一个新的 JWT 实例
		// 解析 token 并获取其中包含的声明
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				// 如果 token 已过期，返回 401 状态码和授权过期提示
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "授权已过期",
				})
				c.Abort() // 中止请求
				return
			}

			// 其他错误返回未登录提示
			c.JSON(http.StatusUnauthorized, "未登陆")
			c.Abort() // 中止请求
			return
		}

		// 将解析出的声明信息存储到上下文中
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next() // 继续处理请求
	}
}

// JWT 是一个结构体，包含签名密钥
type JWT struct {
	SigningKey []byte
}

// 定义一些常见的错误
var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

// 创建一个新的 JWT 实例

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTInfo.SigningKey), // 从全局配置中获取签名密钥
	}
}

//  根据自定义声明创建一个 JWT token

func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	// 使用 HS256 算法生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey) // 使用签名密钥签名并返回 token
}

//  解析 token 并返回其中的声明

func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil // 返回签名密钥进行验证
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed // token 格式错误
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired // token 已过期
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet // token 未激活
			} else {
				return nil, TokenInvalid // token 无效
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil // 返回有效的声明
		}
		return nil, TokenInvalid // token 无效
	} else {
		return nil, TokenInvalid // token 无效
	}
}

// RefreshToken 刷新 token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0) // 将当前时间设置为 Unix 零时间
	}
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil // 返回签名密钥进行验证
	})
	if err != nil {
		return "", err // 返回错误信息
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now                                                // 恢复当前时间
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix() // 更新过期时间
		return j.CreateToken(*claims)                                          // 创建并返回新的 token
	}
	return "", TokenInvalid // 返回无效 token 错误
}
