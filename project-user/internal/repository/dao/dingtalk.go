package dao

import (
	"context"

	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/database"
)

type DingTalkDao struct {
	db *database.GormConn
}

func NewDingTalkDao() *DingTalkDao {
	return &DingTalkDao{
		db: database.NewGormSession(),
	}
}

func (dao *DingTalkDao) InsertDingTalk(conn database.DBConn, ctx context.Context, dingTalk *model.DingTalk) error {
	dao.db = conn.(*database.GormConn)
	return dao.db.Tx(ctx).Create(dingTalk).Error
}
