package database

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"test.com/project-user/config"
)

var db *gorm.DB

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
		zap.L().Error("mysql connect error:", zap.Error(err))
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
