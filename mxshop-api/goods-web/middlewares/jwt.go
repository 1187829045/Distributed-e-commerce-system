package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"mxshop-api/goods-web/global"
	"mxshop-api/goods-web/models"
	"net/http"
	"time"
)

// JWTAuth 返回一个 Gin 中间件，用于验证 JWT 令牌。
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 x-token 字段的值，该字段用于传递 JWT 令牌
		token := c.Request.Header.Get("x-token")
		if token == "" {
			// 如果请求头中没有 x-token 字段，返回未授权的错误响应并中止请求处理
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort()
			return
		}

		// 创建一个新的 JWT 实例
		j := NewJWT()
		// 解析 JWT 令牌包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			// 处理不同类型的 JWT 错误
			if err == TokenExpired {
				// 如果令牌过期，返回授权已过期的错误响应并中止请求处理
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "授权已过期",
				})
				c.Abort()
				return
			}

			// 如果令牌无效或其他错误，返回未登录的错误响应并中止请求处理
			c.JSON(http.StatusUnauthorized, "未登陆")
			c.Abort()
			return
		}

		// 将解析出的 claims 信息和 userId 设置到上下文中，以便后续处理器可以访问
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next() // 继续处理请求
	}
}

// JWT 结构体定义，用于管理 JWT 签名密钥
type JWT struct {
	SigningKey []byte // 签名密钥
}

// 预定义的一些 JWT 错误变量
var (
	TokenExpired     = errors.New("Token is expired")            // 令牌过期错误
	TokenNotValidYet = errors.New("Token not active yet")        // 令牌尚未激活错误
	TokenMalformed   = errors.New("That's not even a token")     // 令牌格式错误
	TokenInvalid     = errors.New("Couldn't handle this token:") // 无法处理令牌错误
)

// NewJWT 创建一个新的 JWT 实例，并设置签名密钥
func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTInfo.SigningKey), // 从全局配置中获取签名密钥，可以设置过期时间
	}
}

// CreateToken 使用指定的声明创建一个 JWT 令牌，并返回签名后的字符串
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	// 使用指定的签名方法和声明创建一个新的 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用签名密钥对令牌进行签名，并返回签名后的字符串
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析给定的 JWT 令牌字符串，返回包含的声明信息
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	// 使用指定的签名密钥解析令牌字符串
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		// 处理 JWT 解析错误
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// 令牌过期
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	// 验证令牌并返回声明信息
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// RefreshToken 刷新给定的 JWT 令牌字符串，返回新的令牌和可能的错误
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	// 设置时间函数为一个过期的时间
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	// 使用指定的签名密钥解析令牌字符串
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	// 如果令牌有效，刷新令牌的过期时间并返回新的令牌
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
