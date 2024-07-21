package dao

import (
	"context"

	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/database"
)

type OrganizationDao struct {
	db *database.GormConn
}

func (dao *OrganizationDao) SelectOrganizationByOrgId(ctx context.Context, pid int64) (bool, error) {
	var count int64
	session := dao.db.Session(ctx)
	err := session.Model(&model.Organization{}).Where("id in ?", pid).Count(&count).Error
	return count > 0, err
}

func NewOrganizationDao() *OrganizationDao {
	return &OrganizationDao{
		db: database.NewGormSession(),
	}
}

func (dao *OrganizationDao) InsertOrganization(conn database.DBConn, ctx context.Context, org *model.Organization) (*model.Organization, error) {
	dao.db = conn.(*database.GormConn)
	err := dao.db.Tx(ctx).Create(&org).Error
	return org, err
}

func (dao *OrganizationDao) SelectOrganizationListByMId(ctx context.Context, mid int64) (org []*model.Organization, err error) {
	result := dao.db.Session(ctx).Where("member_id = ?", mid).Find(&org)
	return org, result.Error
}
