package usersvc

import (
	"context"
	"errors"
	"os/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
)

type UserService struct {
	col *mongo.Collection
}

func NewUserService() *UserService {
	return &UserService{imongo.BizDataBase().Collection("user")}
}

func (s *UserService) FindByLogin(ctx context.Context, username, password string) (*user.User, error) {
	var u user.User
	err := s.col.FindOne(ctx, bson.M{"username": username, "password": password}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
