package database

import (
	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

//func InitRedis(cfg *config.Config) (err error) {
//	Rdb = redis.NewClient(&redis.Options{
//		Addr: fmt.Sprintf("%s:%d",
//			cfg.Redis.Host,
//			cfg.Redis.Port,
//		),
//		Password: cfg.Redis.Password,
//		DB:       cfg.Redis.DB,
//	})
//
//	_, err = Rdb.Ping(context.Background()).Result()
//	return err
//}

func GetRedis() *redis.Client {
	return Rdb
}
