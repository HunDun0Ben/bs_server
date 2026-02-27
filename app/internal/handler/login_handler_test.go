package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/HunDun0Ben/bs_server/app/internal/dto"
	"github.com/HunDun0Ben/bs_server/app/internal/model/user"
	"github.com/HunDun0Ben/bs_server/app/internal/service/authsvc"
	"github.com/HunDun0Ben/bs_server/app/middleware"
	"github.com/HunDun0Ben/bs_server/app/pkg/bsjwt"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"

	confModel "github.com/HunDun0Ben/bs_server/app/pkg/conf/model"
)

func init() {
	// 配置测试所需的 JWT 环境
	conf.AppConfig = confModel.AppConfig{
		JWT: confModel.JWTConfig{
			Secret:        "test-secret-key-1234567890",
			Expire:        time.Hour,
			RefreshExpire: time.Hour * 24,
		},
	}
}

// MockUserService 模拟用户服务逻辑
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FindByLogin(ctx context.Context, username, password string) (*user.User, error) {
	args := m.Called(ctx, username, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserService) EnableMFA(ctx context.Context, username, secret string, recoveryCodes []string) error {
	return m.Called(ctx, username, secret, recoveryCodes).Error(0)
}

func (m *MockUserService) GetMFAInfo(ctx context.Context, username string) (string, bool, error) {
	args := m.Called(ctx, username)
	return args.String(0), args.Bool(1), args.Error(2)
}

func (m *MockUserService) SaveMFASecret(ctx context.Context, username, secret string) error {
	return m.Called(ctx, username, secret).Error(0)
}

func (m *MockUserService) IsHighRisk(ctx context.Context, u *user.User, ip string) (bool, []string) {
	args := m.Called(ctx, u, ip)
	return args.Bool(0), args.Get(1).([]string)
}

func (m *MockUserService) UpdateLoginInfo(ctx context.Context, userID string, ip string) error {
	return m.Called(ctx, userID, ip).Error(0)
}

// MockAuthService 模拟身份验证与令牌管理服务
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) StoreRefreshToken(ctx context.Context, jti, username string, expiration time.Duration) error {
	return m.Called(ctx, jti, username, expiration).Error(0)
}

func (m *MockAuthService) IsRefreshTokenValid(ctx context.Context, jti string) (string, error) {
	args := m.Called(ctx, jti)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) InvalidateRefreshToken(ctx context.Context, jti string) error {
	return m.Called(ctx, jti).Error(0)
}

func (m *MockAuthService) InvalidateAccessToken(ctx context.Context, jti string, remainingTime time.Duration) error {
	return m.Called(ctx, jti, remainingTime).Error(0)
}

func (m *MockAuthService) IsAccessTokenValid(ctx context.Context, jti string) (bool, error) {
	args := m.Called(ctx, jti)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) GenerateTOTPSecret(username string) (*otp.Key, error) {
	args := m.Called(username)
	return args.Get(0).(*otp.Key), args.Error(1)
}

func (m *MockAuthService) GenerateRecoveryCodes() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockAuthService) ValidateTOTPCode(secret string, code string) bool {
	return m.Called(secret, code).Bool(0)
}

func (m *MockAuthService) VerifyAndActivateMFA(ctx context.Context, username, secret, code string) error {
	return m.Called(ctx, username, secret, code).Error(0)
}

func (m *MockAuthService) GetUserMFASecret(ctx context.Context, username string) (string, error) {
	args := m.Called(ctx, username)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) SaveMFASecret(ctx context.Context, username, secret string) error {
	return m.Called(ctx, username, secret).Error(0)
}

// TestLoginHandler_Login_MFA_Required 测试高风险场景下的分阶段认证逻辑。
// 验证当系统识别到风险时，不直接签发正式 Token，而是返回 MFA_REQUIRED 提示。
func TestLoginHandler_Login_MFA_Required(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserSvc := new(MockUserService)
	mockAuthSvc := new(MockAuthService)
	h := NewLoginHandler(mockUserSvc, mockAuthSvc, nil)

	r := gin.New()
	r.Use(middleware.WebErrorHandler())
	r.POST("/login", h.Login)

	testUser := &user.User{
		ID:          "user-123",
		Username:    "testuser",
		LastLoginIP: "1.1.1.1",
	}

	// 模拟核心流程：验证账号 -> 识别风险 -> 生成受限令牌 (此时不更新登录信息)
	mockUserSvc.On("FindByLogin", mock.Anything, "testuser", "password").Return(testUser, nil)
	mockUserSvc.On("IsHighRisk", mock.Anything, testUser, mock.Anything).Return(true, []string{"totp"}).Once()
	// mockUserSvc.On("UpdateLoginInfo", ...) 不应被调用
	mockAuthSvc.On("StoreRefreshToken", mock.Anything, mock.Anything, "testuser", mock.Anything).Return(nil)

	loginReq := dto.LoginRequest{
		Username: "testuser",
		Password: "password",
	}
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 断言：响应应成功返回
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data, ok := resp["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")

	// 断言：业务字段确认需要 MFA
	assert.Equal(t, true, data["mfa_required"])
	assert.Contains(t, data["required_types"], "totp")

	// 断言：验证 Token 内部状态位（mfa_p 为 true）
	tokenStr, ok := data["accessToken"].(string)
	assert.True(t, ok, "accessToken should be a string")
	claims, _ := bsjwt.ParseToken(tokenStr)
	assert.True(t, claims.MFAPending)

	mockUserSvc.AssertExpectations(t)
	mockAuthSvc.AssertExpectations(t)
}

// MockMFAProvider 模拟 MFA 验证器
type MockMFAProvider struct {
	mock.Mock
}

func (m *MockMFAProvider) GetID() string { return m.Called().String(0) }
func (m *MockMFAProvider) Verify(ctx context.Context, secret, code string) (bool, error) {
	args := m.Called(ctx, secret, code)
	return args.Bool(0), args.Error(1)
}

// TestLoginHandler_VerifyMFA_Success 测试 MFA 二次验证通过后的“令牌升级”流程。
// 验证当用户提供正确验证码后，系统作废旧的受限令牌并颁发全新的全权限令牌。
func TestLoginHandler_VerifyMFA_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserSvc := new(MockUserService)
	mockAuthSvc := new(MockAuthService)
	mockProvider := new(MockMFAProvider)
	mockProvider.On("GetID").Return("totp")

	mfaSvc := authsvc.NewMFAVerificationService(mockProvider)
	h := NewLoginHandler(mockUserSvc, mockAuthSvc, mfaSvc)

	// 1. 模拟 context 中已有的受限 Claims (这是中间件校验通过后的状态)
	claims := &bsjwt.CustomClaims{
		Username:      "testuser",
		MFAPending:    true,
		RequiredTypes: []string{"totp"},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ID:        "old-jti",
		},
	}

	r := gin.New()
	r.Use(middleware.WebErrorHandler())
	// 使用自定义中间件注入受限的身份上下文
	r.Use(func(c *gin.Context) {
		c.Set("claims", claims)
		c.Next()
	})
	r.POST("/login/mfa-verify", h.VerifyMFA)

	testUser := &user.User{
		ID:       "user-123",
		Username: "testuser",
	}

	// 模拟验证流程：查找用户 -> 获取 MFA 配置 -> 调用验证器 -> 验证成功 -> 更新登录信息 -> 生成新令牌 -> 作废旧令牌
	mockUserSvc.On("FindByUsername", mock.Anything, "testuser").Return(testUser, nil)
	mockUserSvc.On("GetMFAInfo", mock.Anything, "testuser").Return("secret-123", true, nil)
	mockProvider.On("Verify", mock.Anything, "secret-123", "123456").Return(true, nil)
	mockUserSvc.On("UpdateLoginInfo", mock.Anything, "user-123", mock.Anything).Return(nil)

	mockAuthSvc.On("InvalidateAccessToken", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockAuthSvc.On("StoreRefreshToken", mock.Anything, mock.Anything, "testuser", mock.Anything).Return(nil)

	verifyReq := dto.MFAVerifyRequest{Code: "123456"}
	body, _ := json.Marshal(verifyReq)
	req, _ := http.NewRequest(http.MethodPost, "/login/mfa-verify", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data, ok := resp["data"].(map[string]interface{})
	assert.True(t, ok, "Response data should be a map")

	// 验证颁发了新的令牌
	accessToken, ok := data["accessToken"].(string)
	assert.True(t, ok, "accessToken should be a string")
	newClaims, err := bsjwt.ParseToken(accessToken)
	if err != nil {
		t.Fatalf("Failed to parse access token: %v", err)
	}
	// 断言：新的令牌 mfa_p 为 false，具备完整权限
	assert.False(t, newClaims.MFAPending, "新的 Token 不应再处于 MFAPending 状态")

	mockUserSvc.AssertExpectations(t)
	mockAuthSvc.AssertExpectations(t)
	mockProvider.AssertExpectations(t)
}
