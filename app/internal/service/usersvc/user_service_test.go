package usersvc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/HunDun0Ben/bs_server/app/internal/model/user"
	"github.com/HunDun0Ben/bs_server/app/internal/service/usersvc"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
)

// MockUserRepository is a mock of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) imongo.SingleResult {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(imongo.SingleResult)
}

func (m *MockUserRepository) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

// MockSingleResult is a mock of imongo.SingleResult
type MockSingleResult struct {
	mock.Mock
}

func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func TestFindUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockSR := new(MockSingleResult)
	svc := usersvc.NewUserService(mockRepo)

	ctx := context.Background()
	username := "alice"
	password := "hashed_password_1"

	// Setup expectations
	mockRepo.On("FindOne", ctx, mock.Anything, mock.Anything).Return(mockSR)

	expectedUser := &user.User{Username: username}
	mockSR.On("Decode", mock.Anything).Run(func(args mock.Arguments) {
		u := args.Get(0).(*user.User)
		*u = *expectedUser
	}).Return(nil)

	user, err := svc.FindByLogin(ctx, username, password)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)

	mockRepo.AssertExpectations(t)
	mockSR.AssertExpectations(t)
}

// TestIsHighRisk 验证风险判定逻辑，基于上次登录 IP 与当前请求 IP 的一致性。
func TestIsHighRisk(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := usersvc.NewUserService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name        string
		user        *user.User
		currentIP   string
		expectRisk  bool     // 预期风险判定：true（需要 MFA），false（直接放行）
		expectTypes []string // 预期要求的验证方式
	}{
		{
			name: "First Login (No Last IP)",
			user: &user.User{
				LastLoginIP: "", // 场景：新用户或从未登录过的记录
			},
			currentIP:   "1.1.1.1",
			expectRisk:  false, // 预期：首次登录不强制风险 MFA（视业务场景可调）
			expectTypes: nil,
		},
		{
			name: "Same IP as last time",
			user: &user.User{
				LastLoginIP: "1.1.1.1", // 场景：常用设备、常用网络环境登录
			},
			currentIP:   "1.1.1.1",
			expectRisk:  false, // 预期：低风险，直接颁发正式令牌
			expectTypes: nil,
		},
		{
			name: "Different IP (High Risk)",
			user: &user.User{
				LastLoginIP: "1.1.1.1", // 场景：异地登录、代理 IP 登录、非常用设备登录
			},
			currentIP:   "2.2.2.2",
			expectRisk:  true, // 预期：判定高风险，颁发受限令牌并要求 MFA 验证
			expectTypes: []string{"totp", "sms"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isRisk, types := svc.IsHighRisk(ctx, tt.user, tt.currentIP)
			assert.Equal(t, tt.expectRisk, isRisk)
			assert.Equal(t, tt.expectTypes, types)
		})
	}
}
