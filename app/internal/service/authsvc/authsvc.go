package authsvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/HunDun0Ben/bs_server/app/pkg/data/iredis"
)

const (
	BlockListKeyTemplate    = "blacklist_jti:%s"
	RefreshTokenKeyTemplate = "refresh_token:%s"
)

func StoreRefreshToken(jti, username string, expiration time.Duration) error {
	// 设置 Refresh token 到 redis 中
	redisKey := fmt.Sprintf(RefreshTokenKeyTemplate, jti)
	err := iredis.GetRDB().Set(context.Background(),
		redisKey,
		username,
		expiration).Err()
	return err
}

func IsAcccessTokenValid(jti string) (bool, error) {
	redisKey := fmt.Sprintf(BlockListKeyTemplate, jti)
	err := iredis.GetRDB().Get(context.Background(), redisKey).Err()
	if errors.Is(err, redis.Nil) {
		return true, nil
	} else if err != nil {
		return false, err
	}
	return false, nil
}

func IsRefreshTokenValid(jti string) (string, error) {
	redisKey := fmt.Sprintf(RefreshTokenKeyTemplate, jti)
	storedUsername, err := iredis.GetRDB().Get(context.Background(), redisKey).Result()
	return storedUsername, err
}

// 设置 accessToken jti 阻止登录.
func InvalidateAccessToken(jti string, expiration time.Duration) error {
	redisKey := fmt.Sprintf(BlockListKeyTemplate, jti)
	err := iredis.GetRDB().Set(context.Background(), redisKey, "1", expiration).Err()
	return err
}

func InvalidateRefreshToken(jti string) error {
	redisKey := fmt.Sprintf(RefreshTokenKeyTemplate, jti)
	size, err := iredis.GetRDB().Del(context.Background(), redisKey).Result()
	if err != nil {
		return err
	}
	if size > 0 {
		return nil
	}
	return errors.New("key 不存在")
}
