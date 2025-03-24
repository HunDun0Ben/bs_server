package imongo

import (
	"context"
	"demo/common/conf"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Init() {
	uri := conf.GlobalViper.GetString("mongodb.uri")
	cliOptions := options.
		Client().
		// SetLoggerOptions(options.Logger().SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)).
		ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), cliOptions)
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
