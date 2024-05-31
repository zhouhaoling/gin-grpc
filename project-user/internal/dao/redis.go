package dao

import (
	"context"
	"time"

	"test.com/project-user/config"

	"github.com/go-redis/redis/v8"
)

var RC *RedisCache

func init() {
	rdb := redis.NewClient(config.AppConf.ReadRedisConfig())
	RC = NewRedisCache(rdb)
}

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{
		rdb: rdb,
	}
}

func (r *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	//fmt.Println("r.rdb:", r.rdb)
	err := r.rdb.Set(ctx, key, value, expire).Err()
	return err
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := r.rdb.Get(ctx, key).Result()
	return result, err
}
