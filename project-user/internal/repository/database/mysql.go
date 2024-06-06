package database

import (
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"
	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"test.com/project-user/config"
)

var db *gorm.DB

func init() {
	err := InitMySQL()
	if err != nil {
		log.Fatalln("mysql init error")
	}
}

func InitMySQL() (err error) {
	//配置mysql连接参数
	username := config.AppConf.MySQL.UserName
	password := config.AppConf.MySQL.Password
	host := config.AppConf.MySQL.Host
	port := config.AppConf.MySQL.Port
	dbname := config.AppConf.MySQL.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		zap.L().Error("mysql connect error")
		return err
	}

	//开启数据库连接池
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(config.AppConf.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.AppConf.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifetime)
	return
}

func GetDB() *gorm.DB {
	return db
}

type GormConn struct {
	db *gorm.DB
}

func NewGormConn() *GormConn {
	return &GormConn{
		db: GetDB(),
	}
}

func (g *GormConn) Session(ctx context.Context) *gorm.DB {
	return g.db.Session(&gorm.Session{Context: ctx})
}
