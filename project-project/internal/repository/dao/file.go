package dao

import (
	"context"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository/database"
)

type FileDao struct {
	db *database.GormConn
}

func (dao *FileDao) InsertFileByStruct(ctx context.Context, conn database.DBConn, file *model.File) error {
	dao.db = conn.(*database.GormConn)
	err := dao.db.Tx(ctx).Create(&file).Error
	return err
}

func (dao *FileDao) FindFileByIds(ctx context.Context, ids []int64) (list []*model.File, err error) {
	session := dao.db.Session(ctx)
	err = session.Model(&model.File{}).Where("id in (?)", ids).Find(&list).Error
	return
}

func NewFileDao() *FileDao {
	return &FileDao{db: database.NewGormSession()}
}
