package repo

import (
	"context"

	"test.com/project-project/internal/model"
)

type MenuRepo interface {
	SelectMenus(ctx context.Context) ([]*model.ProjectMenu, error)
}
