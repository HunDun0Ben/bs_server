package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ButterflyRepository interface {
	Collection(name string) *mongo.Collection
	CountDocuments(ctx context.Context, colName string, filter interface{}, opts ...*options.CountOptions) (int64, error)
	Find(ctx context.Context, colName string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, colName string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	InsertOne(ctx context.Context, colName string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, colName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type mongoButterflyRepo struct {
	db *mongo.Database
}

func NewButterflyRepository(db *mongo.Database) ButterflyRepository {
	return &mongoButterflyRepo{db: db}
}

func (r *mongoButterflyRepo) Collection(name string) *mongo.Collection {
	return r.db.Collection(name)
}

func (r *mongoButterflyRepo) CountDocuments(ctx context.Context, colName string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return r.db.Collection(colName).CountDocuments(ctx, filter, opts...)
}

func (r *mongoButterflyRepo) Find(ctx context.Context, colName string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return r.db.Collection(colName).Find(ctx, filter, opts...)
}

func (r *mongoButterflyRepo) FindOne(ctx context.Context, colName string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return r.db.Collection(colName).FindOne(ctx, filter, opts...)
}

func (r *mongoButterflyRepo) InsertOne(ctx context.Context, colName string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return r.db.Collection(colName).InsertOne(ctx, document, opts...)
}

func (r *mongoButterflyRepo) UpdateOne(ctx context.Context, colName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return r.db.Collection(colName).UpdateOne(ctx, filter, update, opts...)
}
