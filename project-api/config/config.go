package config

import (
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
	Etcd  *EtcdConfig
	Jwt   *JwtConfig
}

type JwtConfig struct {
	Secret string
	Issuer string
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
	conf.ReadEtcdConfig()
	conf.ReadJwtConfig()
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

type EtcdConfig struct {
	Addrs []string
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

func (c *Config) ReadJwtConfig() {
	c.Jwt = &JwtConfig{
		Secret: c.viper.GetString("jwt.secret"),
		Issuer: c.viper.GetString("jwt.issuer"),
	}
}
