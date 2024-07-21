package dao

import (
	"context"
	"time"

	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/database"
)

type MemberAccountDao struct {
	db *database.GormConn
}

func (dao *MemberAccountDao) UpdateMemberAccountByLastTime(ctx context.Context, mid int64) error {
	lastTime := time.Now().UnixMilli()
	session := dao.db.Session(ctx)
	err := session.Model(&model.MemberAccount{}).Where("member_code = ?", mid).Update("last_login_time", lastTime).Error
	return err
}

func (dao *MemberAccountDao) InsertMemberAccountByStruct(conn database.DBConn, ctx context.Context, ma *model.MemberAccount) error {
	dao.db = conn.(*database.GormConn)
	err := dao.db.Tx(ctx).Create(&ma).Error
	return err
}

func NewMemberAccountDao() *MemberAccountDao {
	return &MemberAccountDao{
		db: database.NewGormSession(),
	}
}
