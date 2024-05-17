package dao

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var RC *RedisCache

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{
		rdb: rdb,
	}
}

func init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 100,
	})
	RC = NewRedisCache(rdb)
}

func (r *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := r.rdb.Set(ctx, key, value, expire).Err()
	return err
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := r.rdb.Get(ctx, key).Result()
	return result, err
}
