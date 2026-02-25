package repository

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type AuthRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) (int64, error)
}

type redisAuthRepo struct {
	rdb *redis.Client
}

func NewAuthRepository(rdb *redis.Client) AuthRepository {
	return &redisAuthRepo{rdb: rdb}
}

func (r *redisAuthRepo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.rdb.Set(ctx, key, value, expiration).Err()
}

func (r *redisAuthRepo) Get(ctx context.Context, key string) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return val, err
}

func (r *redisAuthRepo) Del(ctx context.Context, key string) (int64, error) {
	return r.rdb.Del(ctx, key).Result()
}
