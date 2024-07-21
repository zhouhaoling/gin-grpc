package service

import (
	"context"

	"github.com/jinzhu/copier"

	"test.com/common/errs"

	"test.com/common/encrypts"

	auth "test.com/project-grpc/auth_grpc"
	"test.com/project-project/internal/domain"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository"
	"test.com/project-project/internal/repository/dao"
	"test.com/project-project/internal/repository/database"
)

type AuthService struct {
	auth.UnimplementedAuthServiceServer
	cache             repo.Cache
	tran              database.Transaction
	projectAuthDomain *domain.ProjectAuthDomain
}

func NewAuthService() *AuthService {
	return &AuthService{
		cache:             dao.NewRedisCache(),
		tran:              repository.NewTransaction(),
		projectAuthDomain: domain.NewProjectAuthDomain(),
	}
}

func (svc *AuthService) AuthList(ctx context.Context, msg *auth.AuthRequest) (*auth.ListAuthMessage, error) {
	oid := encrypts.DecryptInt64(msg.OrganizationCode)
	authList, total, err := svc.projectAuthDomain.AuthListPage(oid, msg.Page, msg.PageSize)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var prList []*auth.AuthProject
	copier.Copy(&prList, authList)
	return &auth.ListAuthMessage{
		List:  prList,
		Total: total,
	}, nil
}

func (svc *AuthService) AuthApply(ctx context.Context, msg *auth.AuthRequest) (*auth.ApplyResponse, error) {
	if msg.Action == "getnode" {
		//获取列表
		list, checkedList, err := svc.projectAuthDomain.AllNodeAndAuth(msg.AuthId)
		if err != nil {
			return nil, errs.GrpcError(err)
		}
		var prList []*auth.ProjectNodeMessages
		copier.Copy(&prList, list)
		return &auth.ApplyResponse{List: prList, CheckedList: checkedList}, nil
	}
	if msg.Action == "save" {
		//保存
		nodes := msg.Nodes
		//先删在存 加事务
		authId := msg.AuthId
		err := svc.tran.Action(func(conn database.DBConn) error {
			err := svc.projectAuthDomain.Save(conn, authId, nodes)
			return err
		})
		if err != nil {
			return nil, errs.GrpcError(err.(*errs.BError))
		}
	}
	return &auth.ApplyResponse{}, nil
}

func (svc *AuthService) AuthNodesByMemberId(ctx context.Context, msg *auth.AuthRequest) (*auth.AuthNodesResponse, error) {
	list, err := svc.projectAuthDomain.AuthNodes(msg.MemberId)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	return &auth.AuthNodesResponse{List: list}, nil
}
