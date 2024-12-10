package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

const (
	secretKey = "mysecretkey12345" // 必须是 16、24 或 32 字节
)

// 加密函数
func encryptURL(originalURL string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// 添加时间戳
	timestamp := time.Now().Unix()
	data := originalURL + "|timestamp=" + strconv.FormatInt(timestamp, 10)

	plainText := []byte(data)
	blockSize := block.BlockSize()
	padding := blockSize - len(plainText)%blockSize
	padText := append(plainText, bytes.Repeat([]byte{byte(padding)}, padding)...)

	cipherText := make([]byte, len(padText))
	mode := cipher.NewCBCEncrypter(block, []byte(secretKey)[:blockSize])
	mode.CryptBlocks(cipherText, padText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// 解密函数
func decryptURL(encryptedURL string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(encryptedURL)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	if len(cipherText)%blockSize != 0 {
		return "", errors.New("invalid ciphertext length")
	}

	plainText := make([]byte, len(cipherText))
	mode := cipher.NewCBCDecrypter(block, []byte(secretKey)[:blockSize])
	mode.CryptBlocks(plainText, cipherText)

	// 去除填充
	padding := int(plainText[len(plainText)-1])
	plainText = plainText[:len(plainText)-padding]

	data := string(plainText)
	parts := strings.Split(data, "|timestamp=")
	if len(parts) != 2 {
		return "", errors.New("invalid encrypted URL format")
	}

	timestamp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", errors.New("invalid timestamp")
	}

	// 校验是否过期
	if time.Now().Unix()-timestamp > 300 { // 5 分钟
		return "", errors.New("URL expired")
	}

	return parts[0], nil
}

func main() {
	r := gin.Default()

	r.GET("/secure", func(c *gin.Context) {
		// 获取完整的请求 URL（包括路径和查询参数）
		originalURL := c.Request.URL.String()
		if originalURL == "" {
			c.JSON(400, gin.H{"error": "url is required"})
			return
		}

		encryptedURL, err := encryptURL(originalURL)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"encrypted_url": encryptedURL})
	})

	r.GET("/validate", func(c *gin.Context) {
		encryptedURL := c.GetHeader("encrypted_url")
		//encryptedURL := c.Query("encrypted_url")
		if encryptedURL == "" {
			c.JSON(400, gin.H{"error": "encrypted_url is required"})
			return
		}

		originalURL, err := decryptURL(encryptedURL)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"original_url": originalURL})
	})

	r.Run(":8080")
}
