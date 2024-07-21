package service

import (
	"context"

	"github.com/jinzhu/copier"
	"test.com/common/encrypts"
	"test.com/common/errs"

	ag "test.com/project-grpc/account_grpc"
	"test.com/project-project/internal/domain"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository"
	"test.com/project-project/internal/repository/dao"
	"test.com/project-project/internal/repository/database"
)

type AccountService struct {
	ag.UnimplementedAccountServiceServer
	cache    repo.Cache
	tran     database.Transaction
	adomain  *domain.AccountDomain
	padomain *domain.ProjectAuthDomain
}

func NewAccountService() *AccountService {
	return &AccountService{
		cache:    dao.NewRedisCache(),
		tran:     repository.NewTransaction(),
		adomain:  domain.NewAccountDomain(),
		padomain: domain.NewProjectAuthDomain(),
	}
}

func (svc *AccountService) Account(ctxs context.Context, msg *ag.AccountRequest) (*ag.AccountResponse, error) {
	//1.去account表查询account
	//2.去auth表查询authList
	accountList, total, err := svc.adomain.AccountList(
		msg.OrganizationCode,
		msg.MemberId,
		msg.Page,
		msg.PageSize,
		msg.DepartmentCode,
		msg.SearchType)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	authList, err := svc.padomain.AuthList(encrypts.DecryptInt64(msg.OrganizationCode))
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var maList []*ag.MemberAccount
	copier.Copy(&maList, accountList)
	var prList []*ag.ProjectAuth
	copier.Copy(&prList, authList)
	return &ag.AccountResponse{
		AccountList: maList,
		AuthList:    prList,
		Total:       total,
	}, nil
}
