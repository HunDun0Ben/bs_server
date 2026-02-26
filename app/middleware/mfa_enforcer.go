package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/pkg/bscxt"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsjwt"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsvo"
)

// MFARequiredError 是一种特殊的错误类型，用于指示需要 MFA
const MFARequiredError = "MFA_REQUIRED"

// MFAEnforcerMiddleware 强制执行 MFA 验证
func MFAEnforcerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Context 中获取 JWT Claims
		claims, exists := c.Get(bscxt.ContextClaimsKey)
		if !exists {
			// 如果没有 Claims，说明 JWTAuthMiddleware 没有正确执行，或者未登录
			c.AbortWithStatusJSON(http.StatusUnauthorized, bsvo.NewAppError(http.StatusUnauthorized, "未提供认证信息", nil, nil))
			return
		}

		customClaims, ok := claims.(*bsjwt.CustomClaims)
		if !ok {
			// Claims 类型不匹配
			c.AbortWithStatusJSON(http.StatusInternalServerError, bsvo.NewAppError(http.StatusInternalServerError, "内部认证错误", nil, nil))
			return
		}

		// 检查 MFAPending 状态
		if customClaims.MFAPending {
			// 白名单路径 (相对于 /api/v1)
			whitelist := map[string]bool{
				"/login/mfa-verify": true,
				"/logout":           true,
			}

			// 获取当前请求的完整路径
			fullPath := c.FullPath()
			// 移除可能存在的版本前缀，例如 /api/v1
			if after, ok0 := strings.CutPrefix(fullPath, "/api/v1"); ok0 {
				fullPath = after
			}

			if !whitelist[fullPath] {
				// 如果不在白名单中，则阻止访问
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":           http.StatusForbidden,
					"message":        "MFA 验证未完成",
					"error":          MFARequiredError,
					"required_types": customClaims.RequiredTypes,
				})
				return
			}
		}

		c.Next()
	}
}
