package dao

import (
	"context"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository/database"
)

type ProjectAuthDao struct {
	db *database.GormConn
}

func (dao *ProjectAuthDao) InsertProjectAuth(ctx context.Context, conn database.DBConn, pa *model.ProjectAuth) error {
	dao.db = conn.(*database.GormConn)
	err := dao.db.Tx(ctx).Create(&pa).Error
	return err
}

func (dao *ProjectAuthDao) FindAuthListPage(ctx context.Context, code int64, page int64, size int64) (list []*model.ProjectAuth, total int64, err error) {
	session := dao.db.Session(ctx)
	err = session.Model(&model.ProjectAuth{}).
		Where("organization_code = ? and status=1", code).
		Limit(int(size)).
		Offset(int((page - 1) * size)).
		Find(&list).Error
	err = session.Model(&model.ProjectAuth{}).
		Where("organization_code = ? and status=1", code).
		Count(&total).Error
	return
}

func (dao *ProjectAuthDao) FindAuthList(ctx context.Context, code int64) (list []*model.ProjectAuth, err error) {
	session := dao.db.Session(ctx)
	err = session.Model(&model.ProjectAuth{}).Where("organization_code = ? and status=1", code).Find(&list).Error
	return
}

func NewProjectAuthDao() *ProjectAuthDao {
	return &ProjectAuthDao{db: database.NewGormSession()}
}
