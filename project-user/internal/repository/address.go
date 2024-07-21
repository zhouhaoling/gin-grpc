package repository

import (
	"context"

	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repo"
	"test.com/project-user/internal/repository/dao"
)

type AddressRepository struct {
	adr repo.AddressRepo
}

func NewAddressRepository() *AddressRepository {
	return &AddressRepository{
		adr: dao.NewAddressDao(),
	}
}

func (repo *AddressRepository) FindAddressByMId(ctx context.Context, mid int64) (model.Address, error) {
	return repo.adr.SelectAddressByMId(ctx, mid)
}
