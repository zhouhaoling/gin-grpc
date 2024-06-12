package database

import (
	"context"
	"fmt"

	pr "test.com/project-project/internal/repository/dao"

	"github.com/go-redis/redis/v8"
	"test.com/project-project/config"
)

func InitRedis(cfg *config.Config) (err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Redis.Host,
			cfg.Redis.Port,
		),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	pr.RC = pr.NewRedisCache(rdb)

	_, err = rdb.Ping(context.Background()).Result()
	return err
}
