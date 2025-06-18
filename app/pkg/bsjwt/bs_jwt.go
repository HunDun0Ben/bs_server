package bsjwt

import (
	"os/user"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token.
func GenerateToken(user user.User) (string, error) {
	claims := Claims{
		user.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.AppConfig.JWT.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// 令牌启用时间
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.AppConfig.JWT.Secret))
}
