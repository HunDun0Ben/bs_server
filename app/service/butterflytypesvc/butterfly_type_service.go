package butterflytypesvc

import (
	"context"

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
