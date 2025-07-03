package utils

import (
	"go-gin/internal/app/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 定义标准声明结构体
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// 生成JWT token的函数
func GenerateToken(userID uint, duration time.Duration) (string, error) {
	expirationTime := time.Now().Add(duration * time.Minute)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),     // 生效时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetGlobalConfig().Service.TokenSecret))
}
