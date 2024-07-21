package repo

import (
	"context"

	"test.com/project-project/internal/model"
)

type DepartmentRepo interface {
	FindDepartmentById(ctx context.Context, id int64) (*model.Department, error)
	ListDepartment(organizationCode int64, parentDepartmentCode int64, page int64, size int64) (list []*model.Department, total int64, err error)
	FindDepartment(ctx context.Context, organizationCode int64, parentDepartmentCode int64, name string) (*model.Department, error)
	Save(dpm *model.Department) error
}
