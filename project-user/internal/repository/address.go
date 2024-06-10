package repository

import "test.com/project-user/internal/repository/dao/mysql"

type AddressRepository struct {
	adr *mysql.AddressDao
}

func NewAddressRepository() *AddressRepository {
	return &AddressRepository{
		adr: mysql.NewAddressDao(),
	}
}
