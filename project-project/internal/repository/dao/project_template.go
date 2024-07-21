package dao

import (
	"context"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository/database"
)

type ProjectTemplateDao struct {
	conn *database.GormConn
}

func (dao *ProjectTemplateDao) InsertProjectTemplate(ctx context.Context, dbConn database.DBConn, list []*model.ProjectTemplate) error {
	dao.conn = dbConn.(*database.GormConn)
	err := dao.conn.Tx(ctx).Model(&model.ProjectTemplate{}).Create(list).Error
	return err
}

func (dao *ProjectTemplateDao) SelectProjectTemplateByMid(ctx context.Context, mid int64) (list []*model.ProjectTemplate, err error) {
	session := dao.conn.Session(ctx)
	err = session.Model(&model.ProjectTemplate{}).Where("member_code = ?", mid).Find(&list).Error
	return
}

var (
	isSystem   = 1
	noIsSystem = 0
)

func (dao *ProjectTemplateDao) SelectProjectTemplateSystem(ctx context.Context, page int64, size int64) (pts []model.ProjectTemplate, total int64, err error) {
	session := dao.conn.Session(ctx)
	index := (page - 1) * size
	err = session.Where("is_system = ?", isSystem).Limit(int(size)).Offset(int(index)).Find(&pts).Error
	session.Model(&model.ProjectTemplate{}).Where("is_system = ?", isSystem).Count(&total)
	return pts, total, err
}

func (dao *ProjectTemplateDao) SelectProjectTemplateCustom(ctx context.Context, mid int64, organizationCode int64, page int64, size int64) (pts []model.ProjectTemplate, total int64, err error) {
	session := dao.conn.Session(ctx)
	index := (page - 1) * size
	err = session.Where("is_system = ? and member_code = ? and organization_code = ?", noIsSystem, mid, organizationCode).
		Limit(int(size)).
		Offset(int(index)).
		Find(&pts).Error

	session.Model(&model.ProjectTemplate{}).Where("is_system = ? and member_code = ? and organization_code = ?", noIsSystem, mid, organizationCode).Count(&total)
	return pts, total, err
}

func (dao *ProjectTemplateDao) SelectProjectTemplateAll(ctx context.Context, organizationCode int64, page int64, size int64) (pts []model.ProjectTemplate, total int64, err error) {
	session := dao.conn.Session(ctx)
	index := (page - 1) * size
	err = session.Where("organization_code = ?", organizationCode).
		Limit(int(size)).
		Offset(int(index)).
		Find(&pts).Error
	session.Model(&model.ProjectTemplate{}).Where("organization_code = ?", organizationCode).Count(&total)
	return pts, total, err
}

func NewProjectTemplateDao() *ProjectTemplateDao {
	return &ProjectTemplateDao{
		conn: database.NewGormSession(),
	}
}
