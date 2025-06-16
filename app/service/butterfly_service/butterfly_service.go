package butterfly_service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/HunDun0Ben/bs_server/app/entities/insect"
	"github.com/HunDun0Ben/bs_server/common/data/imongo"
)

type ButterflyService struct {
	col *mongo.Collection
}

func NewButterflyService() *ButterflyService {
	return &ButterflyService{imongo.BizDataBase().Collection("butterfly_info")}
}

func (s *ButterflyService) GetList(ctx context.Context) ([]insect.Insect, error) {
	cursor, err := s.col.Find(ctx, bson.M{})
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
