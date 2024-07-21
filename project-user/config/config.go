package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"test.com/common/logs"
)

var AppConf = InitConfig()

type Config struct {
	viper *viper.Viper
	App   *AppConfig
	Log   *LogConfig
	Redis *RedisConfig
	Grpc  *GrpcConfig
	Etcd  *EtcdConfig
	MySQL *MySQLConfig
	Snow  *SnowConfig
	Jwt   *JwtConfig
}

type JwtConfig struct {
	Secret     string
	Issuer     string
	TokenType  string
	AccessExp  int
	RefreshExp int
}

type SnowConfig struct {
	StartTime string
	MachineID int64
}

type MySQLConfig struct {
	UserName     string
	Password     string
	Host         string
	DB           string
	Port         int
	MaxOpenConns int
	MaxIdleConns int
}

type GrpcConfig struct {
	Addr    string
	Name    string
	Version string
	Weight  int64
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

type EtcdConfig struct {
	Addrs []string
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
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return nil
	}
	//调用初始化的一些方法
	conf.ReadServerConfig()
	conf.ReadGrpcConfig()
	conf.ReadEtcdConfig()
	conf.ReadMySQLConfig()
	conf.ReadSnowConfig()
	conf.ReadJwtConfig()
	conf.InitZapLog()
	return conf
}

func (c *Config) ReadJwtConfig() {
	c.Jwt = &JwtConfig{
		Secret:     c.viper.GetString("jwt.secret"),
		Issuer:     c.viper.GetString("jwt.issuer"),
		TokenType:  c.viper.GetString("jwt.token_type"),
		AccessExp:  c.viper.GetInt("jwt.access_exp"),
		RefreshExp: c.viper.GetInt("jwt.refresh_exp"),
	}
}

func (c *Config) ReadSnowConfig() {
	c.Snow = &SnowConfig{
		StartTime: c.viper.GetString("snow.start_time"),
		MachineID: c.viper.GetInt64("snow.machine_id"),
	}
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
		Addr:    c.viper.GetString("grpc.addr"),
		Name:    c.viper.GetString("grpc.name"),
		Version: c.viper.GetString("grpc.version"),
		Weight:  c.viper.GetInt64("grpc.weight"),
	}
}

func (c *Config) ReadEtcdConfig() {
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		log.Fatalln(err)
	}
	c.Etcd = &EtcdConfig{
		Addrs: addrs,
	}
}

func (c *Config) ReadMySQLConfig() {
	c.MySQL = &MySQLConfig{
		UserName:     c.viper.GetString("mysql.username"),
		Password:     c.viper.GetString("mysql.password"),
		Host:         c.viper.GetString("mysql.host"),
		DB:           c.viper.GetString("mysql.db"),
		Port:         c.viper.GetInt("mysql.port"),
		MaxOpenConns: c.viper.GetInt("mysql.max_open_conns"),
		MaxIdleConns: c.viper.GetInt("mysql.max_idle_conns"),
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
