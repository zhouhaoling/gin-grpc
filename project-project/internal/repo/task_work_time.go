package repo

import (
	"context"

	"test.com/project-project/internal/model"
)

type TaskWorkTimeRepo interface {
	InsertTaskWorkTime(ctx context.Context, twt *model.TaskWorkTime) error
	SelectTaskWorkTimeList(ctx context.Context, taskCode int64) (list []*model.TaskWorkTime, err error)
}
