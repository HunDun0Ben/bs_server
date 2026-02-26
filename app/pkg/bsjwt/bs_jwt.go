package bsjwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/HunDun0Ben/bs_server/app/internal/model/user"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
)

const (
	AuthHeaderName = "Authorization"
	BearerName     = "Bearer"

	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
)

type CustomClaims struct {
	Username      string   `json:"username"`
	Roles         []string `json:"roles,omitempty"`
	UserID        string   `json:"uid"`
	MFAPending    bool     `json:"mfa_p"`
	RequiredTypes []string `json:"mfa_types"`
	jwt.RegisteredClaims
}

// GenerateAccessToken 只生成 Access Token.
func GenerateAccessToken(user user.User, mfaPending bool, requiredTypes []string) (string, error) {
	claims := CustomClaims{
		user.Username,
		user.Roles,
		user.ID,
		mfaPending,
		requiredTypes,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.AppConfig.JWT.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "access_token",
			ID:        uuid.NewString(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.AppConfig.JWT.Secret))
}

// GenerateRefreshToken 只生成 Refresh Token，不包含任何用户信息.
func GenerateRefreshToken(mfaPending bool, requiredTypes []string) (string, *CustomClaims, error) {
	refreshClaims := &CustomClaims{
		// Username 和 Roles 被有意留空
		Username:      "",
		Roles:         nil,
		MFAPending:    mfaPending,
		RequiredTypes: requiredTypes,
		RegisteredClaims: jwt.RegisteredClaims{
			// Refresh Token 的过期时间更长
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.AppConfig.JWT.RefreshExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   "refresh_token",
			ID:        uuid.NewString(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(conf.AppConfig.JWT.Secret))
	return refreshTokenString, refreshClaims, err
}

// GenerateTokenPair 生成 Access Token 和 Refresh Token.
func GenerateTokenPair(user user.User, mfaPending bool, requiredTypes []string) (accessTokenStr, refreshTokenStr string, claims *CustomClaims, err error) {
	accessToken, err := GenerateAccessToken(user, mfaPending, requiredTypes)
	if err != nil {
		return "", "", nil, err
	}

	refreshToken, claims, err := GenerateRefreshToken(mfaPending, requiredTypes)
	if err != nil {
		return "", "", nil, err
	}
	return accessToken, refreshToken, claims, nil
}

// ParseToken 解析 JWT token.
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(conf.AppConfig.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无效的Token")
}

// ParseTokenByHeader 从请求头中解析和验证 JWT.
func ParseTokenByHeader(authHeader string) (*CustomClaims, error) {
	if authHeader == "" {
		return nil, errors.New("未提供认证信息")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == BearerName) {
		return nil, errors.New("认证格式错误")
	}

	tokenString := parts[1]
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, errors.New("无效的Token")
	}
	return claims, nil
}
