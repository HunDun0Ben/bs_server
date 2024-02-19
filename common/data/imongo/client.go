package imongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Init() {
	uri := `mongodb://localhost:27017`
	client, err := mongo.Connect(context.TODO(),
		options.
			Client().
			// SetLoggerOptions(options.Logger().SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)).
			ApplyURI(uri))
	Client = client
	if err != nil {
		panic(err)
	}
}

func Default() *mongo.Client {
	if Client == nil {
		Init()
	}
	return Client
}

func FileDatabase() *mongo.Database {
	Default()
	return Client.Database("file_db")
}
