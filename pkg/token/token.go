package token

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Config 包括 token 包的配置选项.
type Config struct {
	key         string
	userIdKey   string
	usernameKey string
}

// ErrMissingHeader 表示 `Authorization` 请求头为空.
var ErrMissingHeader = errors.New("the length of the `Authorization` header is zero")

var (
	config = Config{"Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", "userIdKey", "usernameKey"}
	once   sync.Once
)

// Init 设置包级别的配置 config, config 会用于本包后面的 token 签发和解析.
func Init(key string, idKey string, usernameKey string) {
	once.Do(func() {
		if key != "" {
			config.key = key
		}
		if idKey != "" {
			config.userIdKey = idKey
		}
		if usernameKey != "" {
			config.usernameKey = usernameKey
		}
	})
}

// Parse 使用指定的密钥 key 解析 token，解析成功返回 token 上下文，否则报错.
func Parse(tokenString string, key string) (uint, string, error) {
	// 解析 token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// 确保 token 加密算法是预期的加密算法
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(key), nil
	})

	// 解析失败
	if err != nil {
		return 0, "", err
	}

	var userID uint
	var username string
	// 如果解析成功，从token 中取出 token 的主题
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// username 可能存在
		if val, exists := claims[config.usernameKey]; exists {
			if s, ok := val.(string); ok {
				username = s
			}
		}

		// userId 可能不存在（旧 token）
		if val, exists := claims[config.userIdKey]; exists {
			// jwt 中数字会被解析为 float64
			if f, ok := val.(float64); ok {
				userID = uint(f)
			}
		}
	}

	return userID, username, nil
}

// ParseRequest 从请求头中获取令牌，并将其传递给 Parse 函数以解析令牌.
func ParseRequest(c *gin.Context) (uint, string, error) {
	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		return 0, "", ErrMissingHeader
	}

	var t string
	// 从请求头中取出token
	fmt.Sscanf(header, "Bearer %s", &t)

	return Parse(t, config.key)
}

// Sign 使用 jwtSecret 签发 token，token 的 claims 中会存放传入的 subject.
func Sign(userId uint, username string) (tokenString string, err error) {
	// Token 的内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.userIdKey:   userId,
		config.usernameKey: username,
		"nbf":              time.Now().Unix(),
		"iat":              time.Now().Unix(),
		"exp":              time.Now().Add(100000 * time.Hour).Unix(),
	})
	// 签发 token
	tokenString, err = token.SignedString([]byte(config.key))

	return
}
