package domain

import (
	"context"
	"time"

	"test.com/common/errs"
	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository/dao"
)

type TaskDomain struct {
	userRpcDomain *UserRpcDomain
	taskRepo      repo.TaskRepo
}

func (d *TaskDomain) FindProjectIdByTaskId(taskId int64) (int64, bool, *errs.BError) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	task, err := d.taskRepo.SelectTaskById(ctx, taskId)
	if err != nil {
		return 0, false, model.MySQLError
	}
	if task == nil {
		return 0, false, nil
	}
	return task.ProjectCode, true, nil
}

func NewTaskDomain() *TaskDomain {
	return &TaskDomain{
		userRpcDomain: NewUserRpcDomain(),
		taskRepo:      dao.NewTaskDao(),
	}
}
