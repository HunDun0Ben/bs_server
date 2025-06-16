package imongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/HunDun0Ben/bs_server/common/conf"
)

type MongoClient struct {
	client *mongo.Client
	// 保证只初始化一次
	once sync.Once
}

var (
	instance *MongoClient
	mu       sync.Mutex
)

func Client() *MongoClient {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			instance = &MongoClient{}
			instance.Init()
		}
	}
	return instance
}

func (m *MongoClient) Init() error {
	var err error
	m.once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		uri := conf.GlobalViper.GetString("mongodb.uri")
		if uri == "" {
			err = fmt.Errorf("MongoDB URI not configured")
			return
		}
		cliOptions := options.Client().
			SetLoggerOptions(options.Logger().
				SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)).
			ApplyURI(uri)
		m.client, err = mongo.Connect(ctx, cliOptions)
	})
	return err
}

func Database(database string) *mongo.Database {
	return Client().client.Database(database)
}

func BizDataBase() *mongo.Database {
	return Client().client.Database(conf.GlobalViper.GetString("mongodb.biz_db"))
}

func FileDatabase() *mongo.Database {
	return Client().client.Database(conf.GlobalViper.GetString("mongodb.file_db"))
}
