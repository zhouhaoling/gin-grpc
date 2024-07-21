package config

import (
	"fmt"

	"test.com/project-project/internal/repository/database"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var db *gorm.DB

func (c *Config) ReConnMySQL() {
	var err error
	//fmt.Println("c.DbConfig.Separation:", c.DbConfig.Separation)
	if c.DbConfig.Separation {
		//开启读写分离
		username := c.DbConfig.Master.UserName
		password := c.DbConfig.Master.Password
		host := c.DbConfig.Master.Host
		port := c.DbConfig.Master.Port
		dbname := c.DbConfig.Master.DB
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			zap.L().Error("mysql connect error:", zap.Error(err))
			//panic(err)
			return
		}
		//slave
		replicas := []gorm.Dialector{}
		for _, v := range c.DbConfig.Slave {
			username := v.UserName //账号
			password := v.Password //密码
			host := v.Host         //数据库地址，可以是Ip或者域名
			port := v.Port         //数据库端口
			Dbname := v.DB         //数据库名
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
			cfg := mysql.Config{
				DSN: dsn,
			}
			replicas = append(replicas, mysql.New(cfg))
		}
		err = db.Use(dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.New(mysql.Config{
				DSN: dsn,
			})},
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}).
			SetMaxIdleConns(10).
			SetMaxOpenConns(200))
		if err != nil {
			zap.L().Error("启用读写分离失败, error=", zap.Error(err))
			return
			//panic(err)
		}
	} else {
		fmt.Println("不启用读写分离")
		//配置mysql连接参数
		username := c.MySQL.UserName
		password := c.MySQL.Password
		host := c.MySQL.Host
		port := c.MySQL.Port
		dbname := c.MySQL.DB
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		//fmt.Println("dsn", dsn)
		//fmt.Println("db:", db)

		if err != nil {
			zap.L().Error("mysql connect error:", zap.Error(err))
			return
		}

		//开启数据库连接池
		sqlDB, err := db.DB()
		if err != nil {
			zap.L().Error("mysql open DB_pool failed:", zap.Error(err))
			panic(err)
		}
		sqlDB.SetMaxIdleConns(c.MySQL.MaxIdleConns)
		sqlDB.SetMaxOpenConns(c.MySQL.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(MaxLifetime)
	}
	fmt.Println("db:", db)
	database.SetDB(db)
}
