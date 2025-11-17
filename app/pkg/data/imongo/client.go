package imongo

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"

	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
)

var (
	client *mongo.Client
	once   sync.Once
)

// Client 返回 MongoDB 客户端单例。
// 它在首次调用时进行初始化，并且是线程安全的。
func Client() *mongo.Client {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		uri := conf.GlobalViper.GetString("mongodb.uri")
		if uri == "" {
			// 当缺少关键配置时，直接 panic 是合理的，因为程序无法正常运行。
			panic("MongoDB URI not configured")
		}

		cliOptions := options.Client().ApplyURI(uri).SetMonitor(otelmongo.NewMonitor())
		if conf.GlobalViper.GetBool("mongodb.debug") {
			cliOptions.SetLoggerOptions(options.Logger().
				SetComponentLevel(
					options.LogComponentCommand,
					options.LogLevelDebug,
				))
		}

		c, err := mongo.Connect(ctx, cliOptions)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %v", err))
		}

		// Ping 服务器以验证连接是否成功建立。
		if err := c.Ping(ctx, nil); err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v", err))
		}

		slog.Info("Successfully connected to MongoDB")
		client = c
	})
	return client
}

func Database(database string) *mongo.Database {
	return Client().Database(database)
}

func BizDataBase() *mongo.Database {
	return Client().Database(conf.GlobalViper.GetString("mongodb.biz_db"))
}

func FileDatabase() *mongo.Database {
	return Client().Database(conf.GlobalViper.GetString("mongodb.file_db"))
}
