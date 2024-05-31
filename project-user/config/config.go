package config

import (
	"fmt"
	"log"
	"os"

	"test.com/common/logs"

	"github.com/go-redis/redis/v8"

	"github.com/spf13/viper"
)

var AppConf = InitConfig()

type Config struct {
	viper *viper.Viper
	App   *AppConfig
	Log   *LogConfig
	Redis *RedisConfig
	Grpc  *GrpcConfig
}

type GrpcConfig struct {
	Addr string
	Name string
}

type AppConfig struct {
	Name string
	Addr string
}

type LogConfig struct {
	ErrorFileName string
	WarnFileName  string
	InfoFileName  string
	MaxSize       int
	MaxAge        int
	MaxBackups    int
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func InitConfig() *Config {
	v := viper.New()
	conf := &Config{
		viper: v,
	}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("config")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workDir + "/config")

	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	//调用初始化的一些方法
	conf.ReadServerConfig()
	conf.ReadGrpcConfig()
	conf.InitZapLog()
	return conf
}

func (c *Config) ReadServerConfig() {
	c.App = &AppConfig{
		Name: c.viper.GetString("server.name"),
		Addr: c.viper.GetString("server.addr"),
	}

	c.Log = &LogConfig{
		ErrorFileName: c.viper.GetString("zap.error_file_name"),
		WarnFileName:  c.viper.GetString("zap.warn_file_name"),
		InfoFileName:  c.viper.GetString("zap.info_file_name"),
		MaxSize:       c.viper.GetInt("zap.max_size"),
		MaxAge:        c.viper.GetInt("zap.max_age"),
		MaxBackups:    c.viper.GetInt("zap.max_backups"),
	}
	c.Redis = &RedisConfig{
		Host:     c.viper.GetString("redis.host"),
		Port:     c.viper.GetInt("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
	}
}

func (c *Config) ReadGrpcConfig() {
	c.Grpc = &GrpcConfig{
		Addr: c.viper.GetString("grpc.addr"),
		Name: c.viper.GetString("grpc.name"),
	}
}

func (c *Config) ReadRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	}
}

func (c *Config) InitZapLog() {
	//从配置中读取日志配置，初始化日志
	lc := &logs.LogConfig{
		ErrorFileName: c.Log.ErrorFileName,
		InfoFileName:  c.Log.InfoFileName,
		WarnFileName:  c.Log.WarnFileName,
		MaxSize:       c.Log.MaxSize,
		MaxAge:        c.Log.MaxAge,
		MaxBackups:    c.Log.MaxBackups,
	}
	if err := logs.InitLogger(lc); err != nil {
		log.Fatalln("日志初始化失败")
	}
}
