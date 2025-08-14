package imongoutil

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"

	"github.com/HunDun0Ben/bs_server/app/pkg/data/imongo"
)

// UploadFileToGridFS uploads a file to GridFS and returns its ID.
func UploadFileToGridFS(ctx context.Context, fileName string, content []byte) (string, error) {
	db := imongo.FileDatabase()
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return "", fmt.Errorf("failed to create gridfs bucket: %w", err)
	}

	uploadStream, err := bucket.OpenUploadStream(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to open upload stream: %w", err)
	}
	defer uploadStream.Close()

	if _, err := io.Copy(uploadStream, bytes.NewReader(content)); err != nil {
		return "", fmt.Errorf("failed to write to upload stream: %w", err)
	}

	fileID, ok := uploadStream.FileID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert file id to ObjectID, got: %v", uploadStream.FileID)
	}

	return fileID.Hex(), nil
}

// DownloadFileFromGridFS downloads a file from GridFS by its ID.
func DownloadFileFromGridFS(ctx context.Context, fileID primitive.ObjectID) ([]byte, error) {
	db := imongo.FileDatabase()
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create gridfs bucket: %w", err)
	}

	downloadStream, err := bucket.OpenDownloadStream(fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to open download stream: %w", err)
	}
	defer downloadStream.Close()

	content, err := io.ReadAll(downloadStream)
	if err != nil {
		return nil, fmt.Errorf("failed to read from download stream: %w", err)
	}

	return content, nil
}

func StoreFile(cxt context.Context, colName string, file imongo.FileStoreData) (string, error) {
	dbCli := imongo.FileDatabase()
	result, err := dbCli.Collection(colName).InsertOne(cxt, file)
	if err != nil {
		return "", err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("无法将插入的ID转换为ObjectID: %v", result.InsertedID)
}

func GetFile(cxt context.Context, colName, id string) (*imongo.FileStoreData, error) {
	dbCli := imongo.FileDatabase()
	var res imongo.FileStoreData
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}
	err = dbCli.Collection(colName).
		FindOne(cxt, bson.M{"_id": objID}).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
