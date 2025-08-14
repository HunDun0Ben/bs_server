package butterflysvc

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/HunDun0Ben/bs_server/app/internal/model/file"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo/imongoutil"
)

var (
	colButterflyImg        = "butterfly_img"
	colButterflyResizedImg = "butterfly_resized_img"
)

type butterflyImgSvc struct {
	col *mongo.Collection
}

type ButterflyImgSvc interface {
	GetAllList(ctx context.Context, filter any) ([]file.ButterflyFile, error)
	FindOne(ctx context.Context, filter any) (*file.ButterflyFile, error)
}

func NewButterflyImgSvc() ButterflyImgSvc {
	return &butterflyImgSvc{
		col: imongo.FileDatabase().Collection(colButterflyImg),
	}
}

func (s *butterflyImgSvc) GetAllList(ctx context.Context, filter any) ([]file.ButterflyFile, error) {
	return imongoutil.FindAll[file.ButterflyFile](ctx, s.col, filter)
}

func (s *butterflyImgSvc) FindOne(ctx context.Context, filter any) (*file.ButterflyFile, error) {
	result, _ := imongoutil.FindOne[file.ButterflyFile](ctx, s.col, filter)
	return result, nil
}

type butterflyResizedImgSvc struct {
	col *mongo.Collection
}

type ButterflyResizedImgSvc interface {
	Insert(ctx context.Context, file *file.ResizedButteryflyFile) error
	FindOne(ctx context.Context, filter any) (*file.ResizedButteryflyFile, error)
	GetAllList(ctx context.Context, filter any) ([]file.ResizedButteryflyFile, error)
	Update(ctx context.Context, filter any, update any) error
}

func NewButterflyResizedImgSvc() ButterflyResizedImgSvc {
	return &butterflyResizedImgSvc{
		col: imongo.FileDatabase().Collection(colButterflyResizedImg),
	}
}

func (s *butterflyResizedImgSvc) GetAllList(ctx context.Context, filter any) ([]file.ResizedButteryflyFile, error) {
	return imongoutil.FindAll[file.ResizedButteryflyFile](ctx, s.col, filter,
		options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}),
	)
}

func (s *butterflyResizedImgSvc) Insert(ctx context.Context, file *file.ResizedButteryflyFile) error {
	return imongoutil.Insert[butterflyResizedImgSvc](ctx, s.col, file)
}

func (s *butterflyResizedImgSvc) FindOne(ctx context.Context, filter any) (*file.ResizedButteryflyFile, error) {
	result, _ := imongoutil.FindOne[file.ResizedButteryflyFile](ctx, s.col, filter)
	return result, nil
}

func (s *butterflyResizedImgSvc) Update(ctx context.Context, filter, update any) error {
	_, err := s.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
