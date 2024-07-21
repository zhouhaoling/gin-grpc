package database

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

var db *gorm.DB

//func InitMySQL() (err error) {
//	if config.AppConf.DbConfig.Separation {
//		//开启读写分离
//		username := config.AppConf.DbConfig.Master.UserName
//		password := config.AppConf.DbConfig.Master.Password
//		host := config.AppConf.DbConfig.Master.Host
//		port := config.AppConf.DbConfig.Master.Port
//		dbname := config.AppConf.DbConfig.Master.DB
//		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)
//
//		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//			Logger: logger.Default.LogMode(logger.Info),
//			NamingStrategy: schema.NamingStrategy{
//				SingularTable: true,
//			},
//		})
//		if err != nil {
//			zap.L().Error("mysql connect error:", zap.Error(err))
//			return err
//		}
//		//slave
//		replicas := []gorm.Dialector{}
//		for _, v := range config.AppConf.DbConfig.Slave {
//			username := v.UserName //账号
//			password := v.Password //密码
//			host := v.Host         //数据库地址，可以是Ip或者域名
//			port := v.Port         //数据库端口
//			Dbname := v.DB         //数据库名
//			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
//			cfg := mysql.Config{
//				DSN: dsn,
//			}
//			replicas = append(replicas, mysql.New(cfg))
//		}
//		db.Use(dbresolver.Register(dbresolver.Config{
//			Sources: []gorm.Dialector{mysql.New(mysql.Config{
//				DSN: dsn,
//			})},
//			Replicas: replicas,
//			Policy:   dbresolver.RandomPolicy{},
//		}).
//			SetMaxIdleConns(10).
//			SetMaxOpenConns(200))
//	} else {
//		//配置mysql连接参数
//		username := config.AppConf.MySQL.UserName
//		password := config.AppConf.MySQL.Password
//		host := config.AppConf.MySQL.Host
//		port := config.AppConf.MySQL.Port
//		dbname := config.AppConf.MySQL.DB
//		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)
//
//		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//			Logger: logger.Default.LogMode(logger.Info),
//			NamingStrategy: schema.NamingStrategy{
//				SingularTable: true,
//			},
//		})
//		if err != nil {
//			zap.L().Error("mysql connect error:", zap.Error(err))
//			return err
//		}
//
//		//开启数据库连接池
//		sqlDB, err := db.DB()
//		if err != nil {
//			zap.L().Error("mysql open DB_pool failed:", zap.Error(err))
//			return err
//		}
//		sqlDB.SetMaxIdleConns(config.AppConf.MySQL.MaxIdleConns)
//		sqlDB.SetMaxOpenConns(config.AppConf.MySQL.MaxOpenConns)
//		sqlDB.SetConnMaxLifetime(config.MaxLifetime)
//	}
//
//	return
//}

func SetDB(mysqldb *gorm.DB) {
	if mysqldb == nil {
		log.Fatalln(errors.New("错误sql"))
	}
	db = mysqldb
}

func GetDB() *gorm.DB {
	return db
}

// GormConn 开启session db
type GormConn struct {
	gdb *gorm.DB
	tx  *gorm.DB
}

func NewGormSession() *GormConn {
	return &GormConn{
		gdb: GetDB(),
	}
}

// NewTran 开启事务，或者直接使用NewGormSession()然后在赋值
func NewTran() *GormConn {

	return &GormConn{
		gdb: GetDB(),
		tx:  GetDB(),
	}
}

func (g *GormConn) Session(ctx context.Context) *gorm.DB {
	return g.gdb.Session(&gorm.Session{Context: ctx})
}

func (g *GormConn) Rollback() {
	g.tx.Rollback()

}

func (g *GormConn) Commit() {
	g.tx.Commit()

}

// Begin 开启事物,g.gdb.Begin()和GetDB().Begin()有什么区别吗？
func (g *GormConn) Begin() {
	g.tx = g.gdb.Begin()
}

func (g *GormConn) Tx(ctx context.Context) *gorm.DB {
	return g.tx.WithContext(ctx)
}
