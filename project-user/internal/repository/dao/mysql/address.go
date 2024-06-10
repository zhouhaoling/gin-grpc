package mysql

import (
	"context"

	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/database"
)

type AddressDao struct {
	db *database.GormConn
}

func NewAddressDao() *AddressDao {
	return &AddressDao{
		db: database.NewGormSession(),
	}
}

func (dao *AddressDao) InsertAddress(conn database.DBConn, ctx context.Context, address *model.Address) error {
	dao.db = conn.(*database.GormConn)
	return dao.db.Tx(ctx).Create(address).Error
}
