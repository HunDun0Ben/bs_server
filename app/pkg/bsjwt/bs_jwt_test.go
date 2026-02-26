package bsjwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/HunDun0Ben/bs_server/app/internal/model/user"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"

	confModel "github.com/HunDun0Ben/bs_server/app/pkg/conf/model"
)

func init() {
	// 初始化测试配置，确保 JWT 生成所需的 Secret 和过期时间已就绪
	conf.AppConfig = confModel.AppConfig{
		JWT: confModel.JWTConfig{
			Secret:        "test-secret-key-1234567890",
			Expire:        time.Hour,
			RefreshExpire: time.Hour * 24,
		},
	}
}

// TestGenerateAndParseTokenWithMFA 验证 MFA 相关的 Claims 能否正确编码进 Token 并在解析时无损还原
func TestGenerateAndParseTokenWithMFA(t *testing.T) {
	testUser := user.User{
		ID:       "user-123",
		Username: "testuser",
		Roles:    []string{"admin"},
	}

	t.Run("MFA Pending", func(t *testing.T) {
		// 场景：登录触发了 MFA 挂起状态
		mfaPending := true
		requiredTypes := []string{"totp"}

		accessToken, refreshToken, claims, err := GenerateTokenPair(testUser, mfaPending, requiredTypes)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, mfaPending, claims.MFAPending)
		assert.Equal(t, requiredTypes, claims.RequiredTypes)

		// 验证解析 Access Token：确认身份信息与 MFA 状态位一致
		parsedAccess, err := ParseToken(accessToken)
		assert.NoError(t, err)
		assert.Equal(t, testUser.Username, parsedAccess.Username)
		assert.Equal(t, mfaPending, parsedAccess.MFAPending)
		assert.Equal(t, requiredTypes, parsedAccess.RequiredTypes)

		// 验证解析 Refresh Token：Refresh Token 不应携带 Username 敏感信息，但需保留 MFA 状态以供刷新时透传
		parsedRefresh, err := ParseToken(refreshToken)
		assert.NoError(t, err)
		assert.Equal(t, mfaPending, parsedRefresh.MFAPending)
		assert.Equal(t, requiredTypes, parsedRefresh.RequiredTypes)
		assert.Empty(t, parsedRefresh.Username, "Refresh Token 不应包含用户名")
	})

	t.Run("No MFA", func(t *testing.T) {
		// 场景：正常登录，不需要 MFA
		mfaPending := false
		var requiredTypes []string

		accessToken, _, claims, err := GenerateTokenPair(testUser, mfaPending, requiredTypes)
		assert.NoError(t, err)
		assert.False(t, claims.MFAPending)

		parsedAccess, err := ParseToken(accessToken)
		assert.NoError(t, err)
		assert.False(t, parsedAccess.MFAPending)
		assert.Empty(t, parsedAccess.RequiredTypes)
	})
}
