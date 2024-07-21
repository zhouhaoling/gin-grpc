package repo

import (
	"context"

	"test.com/project-project/internal/repository/database"

	"test.com/project-project/internal/model"
)

type FileRepo interface {
	InsertFileByStruct(ctx context.Context, conn database.DBConn, file *model.File) error
	FindFileByIds(ctx context.Context, ids []int64) (list []*model.File, err error)
}
