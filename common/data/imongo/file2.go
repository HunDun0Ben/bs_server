package imongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func StoreFile(cxt context.Context, colName string, file FileStoreData) (bool, error) {
	dbCli := FileDatabase()
	_, err := dbCli.Collection(colName).InsertOne(cxt, file)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetFile(cxt context.Context, colName, id string) (*FileStoreData, error) {
	dbCli := FileDatabase()
	var res FileStoreData
	err := dbCli.Collection(colName).
		FindOne(cxt, bson.M{"_id": id}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
