package usersvc

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
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
