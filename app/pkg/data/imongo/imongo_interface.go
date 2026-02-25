package imongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection 定义了 mongo.Collection 的子集接口，便于 Mock。
type Collection interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
}

// SingleResult 定义了 mongo.SingleResult 的接口。
type SingleResult interface {
	Decode(v interface{}) error
}

// MongoCollection 是 mongo.Collection 的包装器，实现了 Collection 接口。
type MongoCollection struct {
	Col *mongo.Collection
}

func (m *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult {
	return m.Col.FindOne(ctx, filter, opts...)
}

func (m *MongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.Col.UpdateOne(ctx, filter, update, opts...)
}

func (m *MongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return m.Col.InsertOne(ctx, document, opts...)
}

func (m *MongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return m.Col.DeleteOne(ctx, filter, opts...)
}
