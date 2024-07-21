package repo

import (
	"context"

	"test.com/project-project/internal/repository/database"

	"test.com/project-project/internal/model"
)

type ProjectTemplateRepo interface {
	SelectProjectTemplateSystem(ctx context.Context, page int64, size int64) ([]model.ProjectTemplate, int64, error)
	SelectProjectTemplateCustom(ctx context.Context, mid int64, organizationCode int64, page int64, size int64) ([]model.ProjectTemplate, int64, error)
	SelectProjectTemplateAll(ctx context.Context, organizationCode int64, page int64, size int64) ([]model.ProjectTemplate, int64, error)
	SelectProjectTemplateByMid(ctx context.Context, mid int64) ([]*model.ProjectTemplate, error)
	InsertProjectTemplate(ctx context.Context, conn database.DBConn, list []*model.ProjectTemplate) error
}
