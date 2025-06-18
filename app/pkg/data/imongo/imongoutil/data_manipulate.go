package imongoutil

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Count[T any](ctx context.Context, col *mongo.Collection, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return col.CountDocuments(ctx, filter, opts...)
}

func Insert[T any](ctx context.Context, col *mongo.Collection, obj any) error {
	_, err := col.InsertOne(ctx, obj)
	return err
}

// FindAll 是通用的查询函数，T 为目标类型.
func FindAll[T any](ctx context.Context, col *mongo.Collection, filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// FindOne 是通用的查询单个文档的函数，T 为目标类型.
func FindOne[T any](ctx context.Context, col *mongo.Collection, filter interface{}) (*T, error) {
	var result T
	err := col.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // 没有找到文档
		}
		return nil, err // 其他错误
	}
	return &result, nil
}
