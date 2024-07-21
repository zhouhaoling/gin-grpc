package repo

import (
	"context"

	"test.com/project-project/internal/model"
)

type ProjectLogRepo interface {
	SelectProjectLogByTaskCode(ctx context.Context, taskCode int64, comment int) (list []*model.ProjectLog, total int64, err error)
	SelectProjectLogByTaskCodePage(ctx context.Context, taskCode int64, comment int, page int, pageSize int) (list []*model.ProjectLog, total int64, err error)
	InsertProjectLogByStruct(ctx context.Context, pl *model.ProjectLog) error
	SelectProjectLogByMid(ctx context.Context, mid int64, page int64, size int64) (list []*model.ProjectLog, total int64, err error)
}
