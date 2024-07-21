package config

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/spf13/viper"
	"test.com/common/logs"
)

var AppConf = InitConfig()

type Config struct {
	viper    *viper.Viper
	App      *AppConfig
	Log      *LogConfig
	Redis    *RedisConfig
	Grpc     *GrpcConfig
	Etcd     *EtcdConfig
	MySQL    *MySQLConfig
	Snow     *SnowConfig
	Jwt      *JwtConfig
	DbConfig DbConfig
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
	Name         string
}

type DbConfig struct {
	Master     MySQLConfig
	Slave      []MySQLConfig
	Separation bool
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
	//先从nacos读取配置，如果读取不到 再到本地读取
	nacosClient := InitNacosClient()
	configYaml, err := nacosClient.confClient.GetConfig(vo.ConfigParam{
		DataId: "config.yaml",
		Group:  nacosClient.group,
	})
	if err != nil {
		log.Fatalln(err)
	}

	//设置读取文件的类型，默认是JSON类型
	conf.viper.SetConfigType("yaml")
	if configYaml != "" {
		err := conf.viper.ReadConfig(bytes.NewBuffer([]byte(configYaml)))
		if err != nil {
			fmt.Println("读取nacos配置文件失败")
			log.Fatalln(err)
		}
		log.Printf("load nacos config %s \n", configYaml)
		//监听配置文件，发生变动改变配置文件
		err = nacosClient.confClient.ListenConfig(vo.ConfigParam{
			DataId: "config.yaml",
			Group:  nacosClient.group,
			OnChange: func(namespace, group, dataId, data string) {
				log.Printf("load nacos config changed %s \n", data)
				err := conf.viper.ReadConfig(bytes.NewBuffer([]byte(data)))
				if err != nil {
					log.Fatalf("读取修改的nacos配置文件失败 %s", err.Error())
				}
				//所有的配置应该重新读取
				conf.ReLoadAllConfig()
			},
		})
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		//读取本地的配置文件
		workDir, _ := os.Getwd()
		conf.viper.SetConfigName("config")

		conf.viper.AddConfigPath(workDir + "/config")

		err := conf.viper.ReadInConfig()
		if err != nil {
			fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
			panic(err)
		}
	}

	//调用初始化的一些方法
	conf.ReLoadAllConfig()
	return conf
}

// ReLoadAllConfig 重新读取配置
func (c *Config) ReLoadAllConfig() {
	c.ReadServerConfig()
	c.ReadGrpcConfig()
	c.ReadEtcdConfig()
	c.ReadMySQLConfig()
	c.ReadSnowConfig()
	c.ReadJwtConfig()
	c.InitZapLog()
	//c.InitDbConfig()
	//重新创建相关的客户端
	c.ReConnRedis()
	c.ReConnMySQL()
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

func (c *Config) ReadRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
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

func (c *Config) InitDbConfig() {
	mc := DbConfig{}
	mc.Separation = c.viper.GetBool("db.separation")
	var slaves []MySQLConfig
	err := c.viper.UnmarshalKey("db.slave", &slaves)
	if err != nil {
		panic(err)
	}
	master := MySQLConfig{
		UserName: c.viper.GetString("db.master.username"),
		Password: c.viper.GetString("db.master.password"),
		Host:     c.viper.GetString("db.master.host"),
		Port:     c.viper.GetInt("db.master.port"),
		DB:       c.viper.GetString("db.master.db"),
	}
	mc.Master = master
	mc.Slave = slaves
	c.DbConfig = mc
}
