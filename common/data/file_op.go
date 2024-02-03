package data

import (
	"context"
	"demo/common/data/mongodb"
)

func StoreFile(collection string, file FileStoreData) (bool, error) {
	dbCli := mongodb.FileDatabase()
	_, err := dbCli.Collection(collection).InsertOne(context.TODO(), file)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetFile(collection string, file *FileStoreData) {

}
