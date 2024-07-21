package repo

import (
	"context"
	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/database"
)

type AddressRepo interface {
	SelectAddressByMId(ctx context.Context, mid int64) (model.Address, error)
	InsertAddress(conn database.DBConn, ctx context.Context, address *model.Address) error
}
