package domain

import (
	"context"
	"time"

	"go.uber.org/zap"

	"test.com/common/errs"

	"test.com/project-project/internal/repository/dao"

	"test.com/project-project/internal/repo"

	"test.com/project-project/internal/model"
)

type DepartmentDomain struct {
	departmentRepo repo.DepartmentRepo
}

func (d *DepartmentDomain) FindDepartmentById(id int64) (*model.Department, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	dp, err := d.departmentRepo.FindDepartmentById(c, id)
	if err != nil {
		zap.L().Error("DepartmentDomain FindDepartmentById FindDepartmentById() failed", zap.Error(err))
		return nil, model.MySQLError
	}
	return dp, nil
}

func (d *DepartmentDomain) List(organizationCode int64, parentDepartmentCode int64, page int64, size int64) ([]*model.DepartmentDisplay, int64, *errs.BError) {
	list, total, err := d.departmentRepo.ListDepartment(organizationCode, parentDepartmentCode, page, size)
	if err != nil {
		zap.L().Error("DepartmentDomain List ListDepartment() failed", zap.Error(err))
		return nil, 0, model.MySQLError
	}
	var dList []*model.DepartmentDisplay
	for _, v := range list {
		dList = append(dList, v.ToDisplay())
	}
	return dList, total, nil
}

func (d *DepartmentDomain) Save(organizationCode int64, departmentCode int64, parentDepartmentCode int64, name string) (*model.DepartmentDisplay, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	dpm, err := d.departmentRepo.FindDepartment(c, organizationCode, parentDepartmentCode, name)
	if err != nil {
		zap.L().Error("DepartmentDomain Save FindDepartment() failed", zap.Error(err))
		return nil, model.MySQLError
	}
	if dpm == nil {
		dpm = &model.Department{
			Name:             name,
			OrganizationCode: organizationCode,
			CreateTime:       time.Now().UnixMilli(),
		}
		if parentDepartmentCode > 0 {
			dpm.PCode = parentDepartmentCode
		}
		err := d.departmentRepo.Save(dpm)
		if err != nil {
			zap.L().Error("DepartmentDomain Save Save() failed", zap.Error(err))
			return nil, model.MySQLError
		}
		return dpm.ToDisplay(), nil
	}
	return dpm.ToDisplay(), nil
}

func NewDepartmentDomain() *DepartmentDomain {
	return &DepartmentDomain{
		departmentRepo: dao.NewDepartmentDao(),
	}
}
