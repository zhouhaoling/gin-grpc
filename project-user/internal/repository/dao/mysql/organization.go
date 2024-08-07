package mysql

import (
	"context"

	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/database"
)

type OrganizationDao struct {
	db *database.GormConn
}

func NewOrganizationDao() *OrganizationDao {
	return &OrganizationDao{
		db: database.NewGormSession(),
	}
}

func (dao *OrganizationDao) InsertOrganization(conn database.DBConn, ctx context.Context, org *model.Organization) error {
	dao.db = conn.(*database.GormConn)
	return dao.db.Tx(ctx).Create(org).Error
}

func (dao *OrganizationDao) SelectOrganizationListByMId(ctx context.Context, mid int64) (org []*model.Organization, err error) {
	result := dao.db.Session(ctx).Where("member_id = ?", mid).Find(&org)
	return org, result.Error
}
