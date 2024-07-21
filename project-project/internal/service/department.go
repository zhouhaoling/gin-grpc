package service

import (
	"context"

	"test.com/project-project/internal/domain"

	"github.com/jinzhu/copier"
	"test.com/common/encrypts"
	"test.com/common/errs"

	dg "test.com/project-grpc/department_grpc"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository"
	"test.com/project-project/internal/repository/dao"
	"test.com/project-project/internal/repository/database"
)

type DepartmentService struct {
	dg.UnimplementedDepartmentServiceServer
	cache            repo.Cache
	tran             database.Transaction
	departmentDomain *domain.DepartmentDomain
}

func (svc *DepartmentService) Save(ctx context.Context, msg *dg.DepartmentReqMessage) (*dg.DepartmentMessage, error) {
	organizationCode := encrypts.DecryptInt64(msg.OrganizationCode)
	var departmentCode int64
	if msg.DepartmentCode != "" {
		departmentCode = encrypts.DecryptInt64(msg.DepartmentCode)
	}
	var parentDepartmentCode int64
	if msg.ParentDepartmentCode != "" {
		parentDepartmentCode = encrypts.DecryptInt64(msg.ParentDepartmentCode)
	}
	dp, err := svc.departmentDomain.Save(
		organizationCode,
		departmentCode,
		parentDepartmentCode,
		msg.Name)
	if err != nil {
		return &dg.DepartmentMessage{}, errs.GrpcError(err)
	}
	var res = &dg.DepartmentMessage{}
	copier.Copy(res, dp)
	return res, nil
}

func (svc *DepartmentService) List(ctx context.Context, msg *dg.DepartmentReqMessage) (*dg.ListDepartmentMessage, error) {
	organizationCode := encrypts.DecryptInt64(msg.OrganizationCode)
	var parentDepartmentCode int64
	if msg.ParentDepartmentCode != "" {
		parentDepartmentCode = encrypts.DecryptInt64(msg.ParentDepartmentCode)
	}
	dps, total, err := svc.departmentDomain.List(
		organizationCode,
		parentDepartmentCode,
		msg.Page,
		msg.PageSize)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var list []*dg.DepartmentMessage
	copier.Copy(&list, dps)
	return &dg.ListDepartmentMessage{List: list, Total: total}, nil
}

func (svc *DepartmentService) ReadDepartment(ctx context.Context, msg *dg.DepartmentReqMessage) (*dg.DepartmentMessage, error) {
	departmentCode := encrypts.DecryptInt64(msg.DepartmentCode)
	dp, err := svc.departmentDomain.FindDepartmentById(departmentCode)
	if err != nil {
		return &dg.DepartmentMessage{}, errs.GrpcError(err)
	}
	var res = &dg.DepartmentMessage{}
	copier.Copy(res, dp.ToDisplay())
	return res, nil
}

func NewDepartmentService() *DepartmentService {
	return &DepartmentService{
		cache:            dao.NewRedisCache(),
		tran:             repository.NewTransaction(),
		departmentDomain: domain.NewDepartmentDomain(),
	}
}
