package service

import (
	"context"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"test.com/common/errs"

	"test.com/project-project/internal/domain"

	mg "test.com/project-grpc/menu_grpc"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository"
	"test.com/project-project/internal/repository/dao"
	"test.com/project-project/internal/repository/database"
)

type MenuService struct {
	mg.UnimplementedMenuServiceServer
	cache      repo.Cache
	tran       database.Transaction
	menuDomain *domain.MenuDomain
}

func (svc *MenuService) MenuList(ctx context.Context, msg *mg.MenuReqMessage) (*mg.MenuResponseMessage, error) {
	treeList, err := svc.menuDomain.MenuTreeList()
	if err != nil {
		zap.L().Error("MenuList error", zap.Error(err))
		return nil, errs.GrpcError(err)
	}
	var list []*mg.MenusMessage
	copier.Copy(&list, treeList)
	return &mg.MenuResponseMessage{List: list}, nil
}

func NewMenuService() *MenuService {
	return &MenuService{
		cache:      dao.NewRedisCache(),
		tran:       repository.NewTransaction(),
		menuDomain: domain.NewMenuDomain(),
	}
}
