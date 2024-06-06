package database

import (
	"github.com/go-redis/redis/v8"
	"test.com/project-user/config"
	"test.com/project-user/internal/repository/dao/redis"
)

func init() {
	rdb := redis.NewClient(config.AppConf.ReadRedisConfig())
	redis.RC = redis.NewRedisCache(rdb)
}
