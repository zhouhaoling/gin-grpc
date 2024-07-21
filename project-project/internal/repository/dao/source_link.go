package dao

import (
	"context"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository/database"
)

type SourceLinkDao struct {
	db *database.GormConn
}

func (dao *SourceLinkDao) InsertSourceLinkByStruct(ctx context.Context, conn database.DBConn, sl *model.SourceLink) error {
	dao.db = conn.(*database.GormConn)
	err := dao.db.Tx(ctx).Create(&sl).Error
	return err
}
func (dao *SourceLinkDao) SelectSourceLinkByTaskCode(ctx context.Context, taskCode int64) (list []*model.SourceLink, err error) {
	session := dao.db.Session(ctx)
	err = session.Model(&model.SourceLink{}).Where("link_code = ? and link_type = ?", taskCode, "task").Find(&list).Error
	return
}

func NewSourceLinkDao() *SourceLinkDao {
	return &SourceLinkDao{
		db: database.NewGormSession(),
	}
}
