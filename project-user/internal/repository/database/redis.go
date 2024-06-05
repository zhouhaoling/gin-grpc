package database

import (
	"github.com/go-redis/redis/v8"
	"test.com/project-user/config"
	"test.com/project-user/internal/repository/dao/redis_dao"
)

func init() {
	rdb := redis.NewClient(config.AppConf.ReadRedisConfig())
	redis_dao.RC = redis_dao.NewRedisCache(rdb)
}
