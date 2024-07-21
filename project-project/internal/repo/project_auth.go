package repo

import (
	"context"

	"test.com/project-project/internal/repository/database"

	"test.com/project-project/internal/model"
)

type ProjectAuthRepo interface {
	FindAuthList(ctx context.Context, code int64) (list []*model.ProjectAuth, err error)
	FindAuthListPage(ctx context.Context, code int64, page int64, size int64) (list []*model.ProjectAuth, total int64, err error)
	InsertProjectAuth(ctx context.Context, conn database.DBConn, pa *model.ProjectAuth) error
}
