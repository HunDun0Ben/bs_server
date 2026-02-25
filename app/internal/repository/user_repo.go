package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
)

type UserRepository interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) imongo.SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	// Add other necessary methods if needed, or just wrap imongo.Collection
}

type mongoUserRepo struct {
	col imongo.Collection
}

func NewUserRepository(col imongo.Collection) UserRepository {
	return &mongoUserRepo{col: col}
}

func (r *mongoUserRepo) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) imongo.SingleResult {
	return r.col.FindOne(ctx, filter, opts...)
}

func (r *mongoUserRepo) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return r.col.UpdateOne(ctx, filter, update, opts...)
}
