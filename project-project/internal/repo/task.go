package repo

import (
	"context"

	"test.com/project-project/internal/repository/database"

	"test.com/project-project/internal/model"
)

type TaskStagesTemplateRepo interface {
	SelectInProTemIds(ctx context.Context, ids []int) ([]model.MsTaskStagesTemplate, error)
	SelectTaskTemplateByProjectTemplateId(ctx context.Context, tid int) (list []*model.MsTaskStagesTemplate, err error)
	SelectTaskStagesTemplateByPtcodes(ctx context.Context, ptcodes []int64) (list []*model.MsTaskStagesTemplate, err error)
	InsertTaskStagesTemplate(ctx context.Context, conn database.DBConn, list []*model.MsTaskStagesTemplate) error
}

type TaskStagesRepo interface {
	SaveTaskStages(conn database.DBConn, ctx context.Context, ts *model.TaskStages) error
	SelectTaskStagesByProjectId(ctx context.Context, pid, page, pageSize int64) (list []*model.TaskStages, total int64, err error)
	SelectByStageId(ctx context.Context, sid int) (ts *model.TaskStages, err error)
	SelectTaskStagesByPcode(ctx context.Context, code int64) (*model.TaskStages, error)
	InsertTaskStages(ctx context.Context, stage *model.TaskStages) error
	DeleteTaskStagesById(ctx context.Context, code int64) error
	UpdateTaskStagesByStruct(ctx context.Context, code int64, ts *model.TaskStages) error
}

type TaskRepo interface {
	SelectTaskByStageCode(ctx context.Context, stageCode int) ([]*model.Task, error)
	SelectTaskMemberByTaskId(ctx context.Context, taskId int64, mid int64) (task *model.TaskMember, err error)
	SelectTaskMaxIdNum(ctx context.Context, pid int64) (*int, error)
	SelectTaskSortByPIdAndSId(ctx context.Context, pid int64, sid int64) (*int, error)
	InsertTask(conn database.DBConn, ctx context.Context, ts *model.Task) error
	InsertTaskMember(conn database.DBConn, ctx context.Context, tm *model.TaskMember) error
	SelectTaskById(ctx context.Context, taskId int64) (task *model.Task, err error)
	UpdateTaskSort(ctx context.Context, conn database.DBConn, task *model.Task) error
	SelectTaskByStageCodeLtSort(ctx context.Context, stageCode int, sort int) (ts *model.Task, err error)
	SelectTaskByAssignTo(ctx context.Context, mid int64, done int, page int64, pageSize int64) (list []*model.Task, total int64, err error)
	SelectTaskMemberByMemberCode(ctx context.Context, mid int64, done int, page int64, pageSize int64) (list []*model.Task, total int64, err error)
	SelectTaskByCreateBy(ctx context.Context, mid int64, done int, page int64, pageSize int64) (list []*model.Task, total int64, err error)
	SelectTaskMemberByTaskIdAndPage(ctx context.Context, code int64, page int64, pageSize int64) (list []*model.TaskMember, total int64, err error)
	SelectTaskByIds(ctx context.Context, taskIds []int64) (list []*model.Task, err error)
}
