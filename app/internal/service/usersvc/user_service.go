package usersvc

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/HunDun0Ben/bs_server/app/internal/model/user"
	"github.com/HunDun0Ben/bs_server/app/internal/repository"
)

type UserService interface {
	FindByLogin(ctx context.Context, username, password string) (*user.User, error)
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	EnableMFA(ctx context.Context, username, secret string, recoveryCodes []string) error
	GetMFAInfo(ctx context.Context, username string) (string, bool, error)
	SaveMFASecret(ctx context.Context, username, secret string) error
	IsHighRisk(ctx context.Context, u *user.User, ip string) (bool, []string)
	UpdateLoginInfo(ctx context.Context, userID string, ip string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) FindByLogin(ctx context.Context, username, password string) (*user.User, error) {
	var u user.User
	err := s.repo.FindOne(ctx, bson.M{"username": username, "password": password}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (s *userService) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	err := s.repo.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (s *userService) EnableMFA(ctx context.Context, username, secret string, recoveryCodes []string) error {
	_, err := s.repo.UpdateOne(ctx, bson.M{"username": username}, bson.M{
		"$set": bson.M{
			"mfaSecret":     secret,
			"mfaEnabled":    true,
			"recoveryCodes": recoveryCodes,
		},
	})
	return err
}

func (s *userService) GetMFAInfo(ctx context.Context, username string) (string, bool, error) {
	var u user.User
	err := s.repo.FindOne(ctx, bson.M{"username": username}, nil).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", false, nil
		}
		return "", false, err
	}
	return u.MFASecret, u.MFAEnabled, nil
}

func (s *userService) SaveMFASecret(ctx context.Context, username, secret string) error {
	_, err := s.repo.UpdateOne(ctx, bson.M{"username": username}, bson.M{
		"$set": bson.M{
			"mfaSecret":  secret,
			"mfaEnabled": false,
		},
	})
	return err
}

func (s *userService) IsHighRisk(_ context.Context, u *user.User, ip string) (bool, []string) {
	if u.LastLoginIP != "" && u.LastLoginIP != ip {
		return true, []string{"totp", "sms"} // sms is a placeholder
	}
	return false, nil
}

func (s *userService) UpdateLoginInfo(ctx context.Context, userID string, ip string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = s.repo.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": bson.M{
			"lastLoginIP": ip,
			"lastLoginAt": time.Now(),
		},
	})
	return err
}
