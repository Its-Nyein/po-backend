package utilities

import (
	"context"
	"time"

	"po-backend/configs"
)

func RedisSet(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return configs.Envs.Redis.Set(ctx, key, value, expiration).Err()
}

func RedisGet(key string) (string, error) {
	ctx := context.Background()
	return configs.Envs.Redis.Get(ctx, key).Result()
}

func RedisDel(key string) error {
	ctx := context.Background()
	return configs.Envs.Redis.Del(ctx, key).Err()
}
