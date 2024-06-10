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
