package config

import (
	"github.com/go-redis/redis/v8"
	"test.com/project-project/internal/repository/database"
)

func (c *Config) ReConnRedis() {
	rdb := redis.NewClient(c.ReadRedisConfig())
	//rc := &dao.RedisCache{
	//	Rdb: rdb,
	//}
	database.Rdb = rdb
}
