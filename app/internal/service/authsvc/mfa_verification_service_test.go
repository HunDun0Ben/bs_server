package authsvc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMFAProvider 模拟一个具体的 MFA 验证器（如 TOTP 或 SMS 验证器）
type MockMFAProvider struct {
	mock.Mock
}

func (m *MockMFAProvider) GetID() string {
	return m.Called().String(0)
}

func (m *MockMFAProvider) Verify(ctx context.Context, secret string, code string) (bool, error) {
	args := m.Called(ctx, secret, code)
	return args.Bool(0), args.Error(1)
}

// TestMFAVerificationService_VerifyCode 测试 MFAVerificationService 的 Provider 调度和路由逻辑。
// 该服务作为各类型验证器的统一管理层，负责根据 RequiredType 将请求导向正确的实现。
func TestMFAVerificationService_VerifyCode(t *testing.T) {
	mockProvider := new(MockMFAProvider)
	mockProvider.On("GetID").Return("test-mfa")

	// 初始化服务并注入 Mock Provider
	svc := NewMFAVerificationService(mockProvider)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// 场景：验证码正确且 Provider 处理成功
		secret := "secret123"
		code := "123456"
		mockProvider.On("Verify", ctx, secret, code).Return(true, nil).Once()

		valid, err := svc.VerifyCode(ctx, "test-mfa", secret, code)
		assert.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("Failure", func(t *testing.T) {
		// 场景：验证码错误
		secret := "secret123"
		code := "wrong-code"
		mockProvider.On("Verify", ctx, secret, code).Return(false, nil).Once()

		valid, err := svc.VerifyCode(ctx, "test-mfa", secret, code)
		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("Unsupported Provider", func(t *testing.T) {
		// 场景：请求了未在系统中注册的验证器类型
		valid, err := svc.VerifyCode(ctx, "unknown", "s", "c")
		assert.Error(t, err)
		assert.False(t, valid)
		assert.Contains(t, err.Error(), "unsupported MFA provider type")
	})
}
