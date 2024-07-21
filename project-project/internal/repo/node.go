package repo

import (
	"context"

	"test.com/project-project/internal/model"
)

type ProjectNodeRepo interface {
	FindAll(ctx context.Context) (list []*model.ProjectNode, err error)
}
