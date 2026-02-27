package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/internal/service/authsvc"
	"github.com/HunDun0Ben/bs_server/app/pkg/bscxt"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsjwt"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsvo"
)

// JWTAuth JWT 认证中间件.
func JWTAuth(authSvc authsvc.AuthService) gin.HandlerFunc {
	// 解析 JWT
	return func(c *gin.Context) {
		authHeader := c.GetHeader(bsjwt.AuthHeaderName)
		claims, err := bsjwt.ParseTokenByHeader(authHeader)
		if err != nil {
			c.Error(bsvo.NewAppError(http.StatusUnauthorized, err.Error(), nil, err))
			c.Abort()
			return
		}

		is, err := authSvc.IsAccessTokenValid(c, claims.ID)
		if !is {
			c.Error(bsvo.NewAppError(http.StatusUnauthorized, "Token 已失效", nil, nil))
			c.Abort()
			return
		}
		if err != nil {
			c.Error(bsvo.NewAppError(http.StatusInternalServerError, "", nil, err))
			return
		}

		c.Set(bscxt.ContextUsernameKey, claims.Username)
		c.Set(bscxt.ContextRolesKey, claims.Roles)
		c.Set(bscxt.ContextJTIKey, claims.ID)
		c.Set(bscxt.ExpiresAtKey, claims.ExpiresAt.Time)
		c.Set(bscxt.ContextClaimsKey, claims)

		c.Next()
	}
}
