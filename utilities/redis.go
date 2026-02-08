package utilities

import (
	"context"
	"errors"
	"time"

	"po-backend/configs"
)

var errRedisNotConnected = errors.New("redis not connected")

func RedisSet(key string, value interface{}, expiration time.Duration) error {
	if configs.Envs.Redis == nil {
		return errRedisNotConnected
	}
	ctx := context.Background()
	return configs.Envs.Redis.Set(ctx, key, value, expiration).Err()
}

func RedisGet(key string) (string, error) {
	if configs.Envs.Redis == nil {
		return "", errRedisNotConnected
	}
	ctx := context.Background()
	return configs.Envs.Redis.Get(ctx, key).Result()
}

func RedisDel(key string) error {
	if configs.Envs.Redis == nil {
		return errRedisNotConnected
	}
	ctx := context.Background()
	return configs.Envs.Redis.Del(ctx, key).Err()
}
