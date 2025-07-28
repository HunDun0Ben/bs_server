package bsjwt

import (
	"errors"
	"os/user"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateTokenPair 生成 Access Token 和 Refresh Token.
func GenerateTokenPair(user user.User) (map[string]string, error) {
	// 创建 Access Token
	accessClaims := Claims{
		user.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.AppConfig.JWT.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "access_token",
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(conf.AppConfig.JWT.Secret))
	if err != nil {
		return nil, err
	}

	// 创建 Refresh Token
	refreshClaims := Claims{
		user.Username,
		jwt.RegisteredClaims{
			// Refresh Token 的过期时间更长
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.AppConfig.JWT.RefreshExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "refresh_token",
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(conf.AppConfig.JWT.Secret))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
	}, nil
}

// GenerateAccessToken 只生成 Access Token
func GenerateAccessToken(username string) (string, error) {
	claims := Claims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.AppConfig.JWT.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "access_token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.AppConfig.JWT.Secret))
}

// ParseToken 解析 JWT token.
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.AppConfig.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token is invalid")
}
