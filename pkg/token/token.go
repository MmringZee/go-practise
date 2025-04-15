package token

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"sync"
	"time"
)

// Config 包括 token 包的配置选项.
type Config struct {
	// key 用于签发和解析 token 密钥.
	key string
	// identityKey 是 token 中用户身份的键.
	identityKey string
	// expiration 是签发 token 的过期时间.
	// time.Duration 类型, 表示时间段.
	expiration time.Duration
}

// 包内变量
var (
	config = Config{"Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5", "identityKey", 2 * time.Hour}
	once   sync.Once
)

// Init 设置包级别的配置 config, config 会用于本包后面的 token 签发和解析.
func Init(key string, identityKey string, expiration time.Duration) {
	// 在 sync.Once 的作用下, 该闭包只会执行一次
	// 首次调用时会执行闭包内逻辑, 后续所有调用不会再执行闭包
	once.Do(func() {
		if key != "" {
			config.key = key
		}
		if identityKey != "" {
			config.identityKey = identityKey
		}
		if expiration != 0 {
			config.expiration = expiration
		}
	})
}

// Parse 使用指定的密钥 key 解析 token, 成功则返回 token 身份键; 否则报错.
func Parse(tokenString string, key string) (string, error) {
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保 token 加密算法是预期的算法 (断言判断)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		// 返回密钥
		return []byte(key), nil
	})

	// 解析失败
	if err != nil {
		return "", err
	}

	// 解析成功, 则从 token 中取出 token 的主题
	var identityKey string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if key, exists := claims[config.identityKey]; exists {
			if identity, valid := key.(string); valid {
				// 获得身份键
				identityKey = identity
			}
		}
	}
	if identityKey == "" {
		return "", jwt.ErrSignatureInvalid
	}

	return identityKey, nil
}

// ParseRequest 从请求头获取 token, 并传递给 Parse 函数以解析令牌
func ParseRequest(c *gin.Context) (string, error) {
	// 从头部获取 token (一般 token 存放在 "Authorization")
	header := c.Request.Header.Get("Authorization")

	// HTTP 请求中若没有字段 "Authorization", 则 header 会获得空字符串
	if len(header) == 0 {
		return "", errors.New("the length of the `Authorization` head is zero")
	}

	var token string
	// 从请求头取出 token
	_, err := fmt.Sscanf(header, "Bearer %s", &token)
	if err != nil {
		return "", errors.New("the Authorization token cannot be parsed into the specified structure")
	}

	return Parse(token, config.key)
}

// Sign 使用 jwtSecret 签发 token，token 的 claims 中会存放传入的 subject.
func Sign(identityKey string) (string, time.Time, error) {
	// 计算过期时间
	expireAt := time.Now().Add(config.expiration)

	// Token 内容
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,       // 存放用户身份
		"nbf":              time.Now().Unix(), // token 生效时间
		"iat":              time.Now().Unix(), // token 签发时间
		"exp":              time.Now().Unix(), // token 过期时间
	})
	if config.key == "" {
		return "", time.Time{}, jwt.ErrInvalidKey
	}

	// 签发 token
	tokenString, err := token.SignedString([]byte(config.key))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expireAt, nil
}
