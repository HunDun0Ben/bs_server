package imongoutil

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
)

func StoreFile(cxt context.Context, colName string, file imongo.FileStoreData) (bool, error) {
	dbCli := imongo.FileDatabase()
	_, err := dbCli.Collection(colName).InsertOne(cxt, file)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetFile(cxt context.Context, colName, id string) (*imongo.FileStoreData, error) {
	dbCli := imongo.FileDatabase()
	var res imongo.FileStoreData
	err := dbCli.Collection(colName).
		FindOne(cxt, bson.M{"_id": id}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
