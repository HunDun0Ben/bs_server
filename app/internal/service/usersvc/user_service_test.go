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
