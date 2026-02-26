package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/HunDun0Ben/bs_server/app/pkg/bscxt"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsjwt"
)

// TestMFAEnforcerMiddleware 测试 MFA 强制执行中间件在不同用户状态和访问路径下的拦截/放行策略
func TestMFAEnforcerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mfaPending     bool   // Token 中的 MFA 挂起状态
		path           string // 访问路径
		expectedStatus int    // 预期响应码
	}{
		{
			name:           "MFAPending=false, Access Business API",
			mfaPending:     false,
			path:           "/api/v1/user/info",
			expectedStatus: http.StatusOK, // 正常状态，允许访问业务接口
		},
		{
			name:           "MFAPending=true, Access Business API (Blocked)",
			mfaPending:     true,
			path:           "/api/v1/user/info",
			expectedStatus: http.StatusForbidden, // MFA 挂起，拦截业务接口
		},
		{
			name:           "MFAPending=true, Access Whitelist (MFA Verify)",
			mfaPending:     true,
			path:           "/api/v1/login/mfa-verify",
			expectedStatus: http.StatusOK, // 白名单路径：允许进行二次验证
		},
		{
			name:           "MFAPending=true, Access Whitelist (Logout)",
			mfaPending:     true,
			path:           "/api/v1/logout",
			expectedStatus: http.StatusOK, // 白名单路径：允许未验证用户登出
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			// 模拟前置中间件（如 JWTAuth）将 Claims 注入 Context
			r.Use(func(c *gin.Context) {
				c.Set(bscxt.ContextClaimsKey, &bsjwt.CustomClaims{
					MFAPending:    tt.mfaPending,
					RequiredTypes: []string{"totp"},
				})
				c.Next()
			})
			r.Use(MFAEnforcerMiddleware())

			// 模拟处理器，如果请求通过中间件，则返回 200
			handler := func(c *gin.Context) { c.Status(http.StatusOK) }
			r.GET(tt.path, handler)
			r.POST(tt.path, handler)

			// 针对某些白名单路径使用 POST 请求
			method := http.MethodGet
			if tt.path == "/api/v1/login/mfa-verify" || tt.path == "/api/v1/logout" {
				method = http.MethodPost
			}

			req, _ := http.NewRequest(method, tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
