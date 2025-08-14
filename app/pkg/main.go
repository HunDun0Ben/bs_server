package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// uri := os.Getenv("MONGODB_URI")
	uri := `mongodb://localhost:27017`
	if uri == "" {
		slog.Error("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("file_db").Collection("butterfly_img")
	title := "Back to the Future"

	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{Key: "file_name", Value: "0010001.png"}}).
		Decode(&result)
	if err == mongo.ErrNoDocuments {
		slog.Error("No document was found with the ", "title", title)
		return
	}
	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	slog.Info("", slog.Attr{
		Key:   "jsonData",
		Value: slog.StringValue(string(jsonData)),
	})
}

func logger() {
	loggerOptions := options.
		Logger().
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	uri := `mongodb://localhost:27017`

	clientOptions := options.
		Client().
		ApplyURI(uri).
		SetLoggerOptions(loggerOptions)

	mongo.Connect(context.TODO(), clientOptions)
}
