package iredis

import (
	"context"
	"testing"
)

// TestGetRDB 测试 GetRDB 函数是否能成功返回一个可用的 Redis 客户端实例。
// 注意：此测试需要一个正在运行的 Redis 实例，并依赖于正确的配置文件 (conf/redis.yaml)。
func TestGetRDB(t *testing.T) {
	// 第一次调用以初始化
	rdb := GetRDB()
	if rdb == nil {
		t.Fatal("GetRDB() returned nil, expected a redis client instance")
	}

	// 再次调用以测试单例模式
	rdb2 := GetRDB()
	if rdb != rdb2 {
		t.Fatal("GetRDB() should return a singleton instance, but got different instances")
	}

	// 测试与 Redis 服务器的连通性
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		t.Fatalf("Failed to ping Redis server: %v", err)
	}

	if pong != "PONG" {
		t.Errorf("Expected PONG from Redis, but got '%s'", pong)
	}

	t.Log("Successfully connected to Redis and received PONG.")
}
