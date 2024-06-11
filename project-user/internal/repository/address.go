package repository

import (
	"context"

	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/dao/mysql"
)

type AddressRepository struct {
	adr *mysql.AddressDao
}

func NewAddressRepository() *AddressRepository {
	return &AddressRepository{
		adr: mysql.NewAddressDao(),
	}
}

func (repo *AddressRepository) FindAddressByMId(ctx context.Context, mid int64) (model.Address, error) {
	return repo.adr.SelectAddressByMId(ctx, mid)
}
