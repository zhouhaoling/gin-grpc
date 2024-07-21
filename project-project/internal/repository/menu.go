package repository

import (
	"context"
	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository/dao"
)

type MenuRepository struct {
	mdao repo.MenuRepo
}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{
		mdao: dao.NewMenuDao(),
	}
}

func (repo *MenuRepository) FindMenu(ctx context.Context) ([]*model.ProjectMenu, error) {
	return repo.mdao.SelectMenus(ctx)
}
