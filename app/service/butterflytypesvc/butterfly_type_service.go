package butterflytypesvc

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/HunDun0Ben/bs_server/app/entities/insect"
	"github.com/HunDun0Ben/bs_server/common/data/imongo"
	"github.com/HunDun0Ben/bs_server/common/data/imongo/imongoutil"
)

const colName = "butterfly_type_info"

type ButterflyTypeService struct {
	col *mongo.Collection
}

func NewButterflyService() *ButterflyTypeService {
	return &ButterflyTypeService{imongo.BizDataBase().Collection(colName)}
}

func (s *ButterflyTypeService) GetAllList(ctx context.Context) ([]insect.Insect, error) {
	return imongoutil.FindAll[insect.Insect](ctx, s.col, bson.M{})
}

func (s *ButterflyTypeService) InitAll(ctx context.Context, list []insect.Insect) error {
	if len(list) == 0 {
		return nil
	}

	ct, err := s.col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if ct > 0 {
		return errors.New("数据表已存在数据，无法初始化")
	}

	docs := make([]interface{}, len(list))
	for i, v := range list {
		docs[i] = v
	}
	if _, err := s.col.InsertMany(context.Background(), docs); err != nil {
		return err
	}
	return nil
}

func (s *ButterflyTypeService) Update(cxt context.Context) {
	return
}
