package database

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"test.com/project-user/config"
	myProjectRedis "test.com/project-user/internal/repository/dao/redis_dao"
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
	myProjectRedis.RC = myProjectRedis.NewRedisCache(rdb)

	_, err = rdb.Ping(context.Background()).Result()
	return err
}
