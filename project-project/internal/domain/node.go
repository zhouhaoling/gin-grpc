package domain

import (
	"context"

	"go.uber.org/zap"

	"test.com/project-api/config"
	"test.com/project-project/internal/repository/dao"

	"test.com/project-project/internal/repo"

	"test.com/project-project/internal/model"

	"test.com/common/errs"
)

type ProjectNodeDomain struct {
	projectNodeRepo repo.ProjectNodeRepo
}

func (d *ProjectNodeDomain) TreeList() ([]*model.ProjectNodeTree, *errs.BError) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeOut)
	defer cancel()
	nodes, err := d.projectNodeRepo.FindAll(ctx)
	if err != nil {
		return nil, model.MySQLError
	}
	treeList := model.ToNodeTreeList(nodes)
	return treeList, nil
}

func (d *ProjectNodeDomain) AllNodeList() ([]*model.ProjectNode, *errs.BError) {
	nodes, err := d.projectNodeRepo.FindAll(context.Background())
	if err != nil {
		zap.L().Error("project AllNodeList projectNodeRepo.FindAll() error", zap.Error(err))
		return nil, model.MySQLError
	}
	return nodes, nil
}

func NewProjectNodeDomain() *ProjectNodeDomain {
	return &ProjectNodeDomain{
		projectNodeRepo: dao.NewProjectNodeDao(),
	}
}
