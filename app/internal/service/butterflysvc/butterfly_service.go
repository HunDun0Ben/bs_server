package butterflysvc

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/HunDun0Ben/bs_server/app/internal/model/file"
	"github.com/HunDun0Ben/bs_server/app/internal/model/insect"
	"github.com/HunDun0Ben/bs_server/app/internal/repository"
)

const (
	colButterflyTypeInfo   = "butterfly_type_info"
	colButterflyImg        = "butterfly_img"
	colButterflyResizedImg = "butterfly_resized_img"
)

type ButterflyService interface {
	// Type Info
	CountTypes(ctx context.Context) (int64, error)
	GetTypes(ctx context.Context) ([]insect.Insect, error)
	InitTypes(ctx context.Context, list []insect.Insect) error

	// Img Info
	GetImgs(ctx context.Context, filter any) ([]file.ButterflyFile, error)
	FindImg(ctx context.Context, filter any) (*file.ButterflyFile, error)

	// Resized Img
	InsertResizedImg(ctx context.Context, file *file.ResizedButteryflyFile) error
	FindResizedImg(ctx context.Context, filter any) (*file.ResizedButteryflyFile, error)
	GetResizedImgs(ctx context.Context, filter any) ([]file.ResizedButteryflyFile, error)
	UpdateResizedImg(ctx context.Context, filter any, update any) error
}

type butterflyService struct {
	repo repository.ButterflyRepository
}

func NewButterflyService(repo repository.ButterflyRepository) ButterflyService {
	return &butterflyService{repo: repo}
}

func (s *butterflyService) CountTypes(ctx context.Context) (int64, error) {
	return s.repo.CountDocuments(ctx, colButterflyTypeInfo, bson.M{})
}

func (s *butterflyService) GetTypes(ctx context.Context) ([]insect.Insect, error) {
	cursor, err := s.repo.Find(ctx, colButterflyTypeInfo, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []insect.Insect
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *butterflyService) InitTypes(ctx context.Context, list []insect.Insect) error {
	if len(list) == 0 {
		return nil
	}
	ct, err := s.repo.CountDocuments(ctx, colButterflyTypeInfo, bson.M{})
	if err != nil {
		return err
	}
	if ct > 0 {
		return errors.New("数据表已存在数据，无法初始化")
	}
	docs := make([]any, len(list))
	for i, v := range list {
		docs[i] = v
	}
	_, err = s.repo.Collection(colButterflyTypeInfo).InsertMany(ctx, docs)
	return err
}

func (s *butterflyService) GetImgs(ctx context.Context, filter any) ([]file.ButterflyFile, error) {
	cursor, err := s.repo.Find(ctx, colButterflyImg, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []file.ButterflyFile
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *butterflyService) FindImg(ctx context.Context, filter any) (*file.ButterflyFile, error) {
	var result file.ButterflyFile
	err := s.repo.FindOne(ctx, colButterflyImg, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *butterflyService) InsertResizedImg(ctx context.Context, file *file.ResizedButteryflyFile) error {
	_, err := s.repo.InsertOne(ctx, colButterflyResizedImg, file)
	return err
}

func (s *butterflyService) FindResizedImg(ctx context.Context, filter any) (*file.ResizedButteryflyFile, error) {
	var result file.ResizedButteryflyFile
	err := s.repo.FindOne(ctx, colButterflyResizedImg, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *butterflyService) GetResizedImgs(ctx context.Context, filter any) ([]file.ResizedButteryflyFile, error) {
	cursor, err := s.repo.Find(ctx, colButterflyResizedImg, filter, options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []file.ResizedButteryflyFile
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *butterflyService) UpdateResizedImg(ctx context.Context, filter any, update any) error {
	_, err := s.repo.UpdateOne(ctx, colButterflyResizedImg, filter, update)
	return err
}
