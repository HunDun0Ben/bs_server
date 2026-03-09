package iredis

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
)

var (
	rdb  *redis.Client
	once sync.Once
)

// GetRDB 返回 Redis 客户端单例.
func GetRDB() *redis.Client {
	once.Do(func() {
		cfg := conf.AppConfig.Redis
		rdb = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
			PoolSize: cfg.PoolSize,
		})

		// Enable tracing instrumentation.
		if err := redisotel.InstrumentTracing(rdb); err != nil {
			panic(err)
		}

		// Enable metrics instrumentation.
		if err := redisotel.InstrumentMetrics(rdb); err != nil {
			panic(err)
		}

		// 检查连接
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			panic(fmt.Sprintf("无法连接到 Redis: %v", err))
		}
	})
	return rdb
}
