package domain

import (
	"context"

	"test.com/project-project/internal/repository/database"

	"go.uber.org/zap"

	"test.com/common/errs"
	"test.com/project-project/internal/model"

	"test.com/project-api/config"

	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository/dao"
)

type ProjectAuthNodeDomain struct {
	projectAuthNodeRepo repo.ProjectAuthNodeRepo
}

func (d *ProjectAuthNodeDomain) AuthNodeList(authId int64) ([]string, *errs.BError) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeOut)
	defer cancel()
	list, err := d.projectAuthNodeRepo.FindNodeStringList(ctx, authId)
	if err != nil {
		zap.L().Error("project AuthNodeList ProjectAuthNodeDomain.AuthNodeList() error", zap.Error(err))
		return nil, model.MySQLError
	}
	return list, nil
}

func (d *ProjectAuthNodeDomain) Save(ctx context.Context, conn database.DBConn, authId int64, nodes []string) *errs.BError {
	err := d.projectAuthNodeRepo.DeleteByAuthId(ctx, conn, authId)
	if err != nil {
		zap.L().Error("project Save ProjectAuthNodeDomain.DeleteByAuthId() error", zap.Error(err))
		return model.MySQLError
	}
	err = d.projectAuthNodeRepo.Save(ctx, conn, authId, nodes)
	if err != nil {
		zap.L().Error("project Save ProjectAuthNodeDomain.Save() error", zap.Error(err))
		return model.MySQLError
	}
	return nil
}

func (d *ProjectAuthNodeDomain) FindProjectAuthNode(ctx context.Context) ([]*model.ProjectAuthNode, error) {
	projectAuthNodeList, err := d.projectAuthNodeRepo.SelectProjectAuthNode(ctx)
	if err != nil {
		zap.L().Error("project FindProjectAuthNode ProjectAuthNodeDomain.SelectProjectAuthNode() error", zap.Error(err))
		return nil, err
	}
	return projectAuthNodeList, nil
}

func (d *ProjectAuthNodeDomain) CreateProjectAuthNode(ctx context.Context, conn database.DBConn, auth *model.ProjectAuth, list []*model.ProjectAuthNode) error {
	newProjectAuthNodeList := model.NewProjectAuthNode(list, auth.Id)
	err := d.projectAuthNodeRepo.InsertProjectAuthNodeList(ctx, conn, newProjectAuthNodeList)
	if err != nil {
		zap.L().Error("project CreateProjectAuthNode ProjectAuthNodeDomain.InsertProjectAuthNodeList() error", zap.Error(err))
		return err
	}
	return nil
}

func NewProjectAuthNodeDomain() *ProjectAuthNodeDomain {
	return &ProjectAuthNodeDomain{
		projectAuthNodeRepo: dao.NewProjectAuthNodeDao(),
	}
}
