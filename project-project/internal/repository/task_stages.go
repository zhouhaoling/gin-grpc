package repository

import (
	"context"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository/dao"
	"test.com/project-project/internal/repository/database"
)

type TaskStagesRepository struct {
	tsr repo.TaskStagesRepo
	tr  repo.TaskRepo
	twt repo.TaskWorkTimeRepo
}

func NewTaskStagesRepository() *TaskStagesRepository {
	return &TaskStagesRepository{
		tsr: dao.NewTaskStagesDao(),
		tr:  dao.NewTaskDao(),
		twt: dao.NewTaskWorkTimeDao(),
	}
}

func (repo *TaskStagesRepository) CreateTaskStage(conn database.DBConn, ctx context.Context, ts *model.TaskStages) error {
	return repo.tsr.SaveTaskStages(conn, ctx, ts)
}

func (repo *TaskStagesRepository) FindTaskStagesByProjectId(ctx context.Context, pid, page, pageSize int64) ([]*model.TaskStages, int64, error) {
	return repo.tsr.SelectTaskStagesByProjectId(ctx, pid, page, pageSize)
}

func (repo *TaskStagesRepository) FindTaskByStageCode(ctx context.Context, stageCode int) ([]*model.Task, error) {
	return repo.tr.SelectTaskByStageCode(ctx, stageCode)
}

func (repo *TaskStagesRepository) FindTaskMemberByTaskId(ctx context.Context, taskId int64, mid int64) (*model.TaskMember, error) {
	return repo.tr.SelectTaskMemberByTaskId(ctx, taskId, mid)
}

func (repo *TaskStagesRepository) FindById(ctx context.Context, sid int) (*model.TaskStages, error) {
	return repo.tsr.SelectByStageId(ctx, sid)
}

func (repo *TaskStagesRepository) FindTaskMaxIdNum(ctx context.Context, pid int64) (*int, error) {
	return repo.tr.SelectTaskMaxIdNum(ctx, pid)
}

func (repo *TaskStagesRepository) FindTaskSort(ctx context.Context, pid int64, sid int64) (*int, error) {
	return repo.tr.SelectTaskSortByPIdAndSId(ctx, pid, sid)
}

func (repo *TaskStagesRepository) SaveTask(ctx context.Context, conn database.DBConn, ts *model.Task) error {
	return repo.tr.InsertTask(conn, ctx, ts)
}

func (repo *TaskStagesRepository) SaveTaskMember(ctx context.Context, conn database.DBConn, tm *model.TaskMember) error {
	return repo.tr.InsertTaskMember(conn, ctx, tm)
}

func (repo *TaskStagesRepository) FindTaskById(ctx context.Context, taskId int64) (ts *model.Task, err error) {
	return repo.tr.SelectTaskById(ctx, taskId)
}

func (repo *TaskStagesRepository) UpdateTaskSort(ctx context.Context, conn database.DBConn, task *model.Task) error {
	return repo.tr.UpdateTaskSort(ctx, conn, task)
}

func (repo *TaskStagesRepository) FindTaskByStageCodeLtSort(ctx context.Context, stageCode int, sort int) (ts *model.Task, err error) {
	return repo.tr.SelectTaskByStageCodeLtSort(ctx, stageCode, sort)
}

func (repo *TaskStagesRepository) FindTaskByAssignTo(ctx context.Context, mid int64, done int, page int64, pageSize int64) ([]*model.Task, int64, error) {
	return repo.tr.SelectTaskByAssignTo(ctx, mid, done, page, pageSize)
}

func (repo *TaskStagesRepository) FindTaskByMemberCode(ctx context.Context, mid int64, done int, page int64, pageSize int64) ([]*model.Task, int64, error) {
	return repo.tr.SelectTaskMemberByMemberCode(ctx, mid, done, page, pageSize)
}
func (repo *TaskStagesRepository) FindTaskByCreateBy(ctx context.Context, mid int64, done int, page int64, pageSize int64) ([]*model.Task, int64, error) {
	return repo.tr.SelectTaskByCreateBy(ctx, mid, done, page, pageSize)
}

func (repo *TaskStagesRepository) FindTaskMemberByTaskIdAndPage(ctx context.Context, code int64, page int64, size int64) (list []*model.TaskMember, total int64, err error) {
	return repo.tr.SelectTaskMemberByTaskIdAndPage(ctx, code, page, size)
}

func (repo *TaskStagesRepository) FindWorkTimeList(ctx context.Context, code int64) ([]*model.TaskWorkTime, error) {
	return repo.twt.SelectTaskWorkTimeList(ctx, code)
}

func (repo *TaskStagesRepository) SaveTaskWorkTime(ctx context.Context, tmt *model.TaskWorkTime) error {
	return repo.twt.InsertTaskWorkTime(ctx, tmt)
}

func (repo *TaskStagesRepository) FindTaskByIds(ctx context.Context, list []int64) ([]*model.Task, error) {
	return repo.tr.SelectTaskByIds(ctx, list)
}

func (repo *TaskStagesRepository) FindTaskStagesByPcode(ctx context.Context, code int64) (*model.TaskStages, error) {
	return repo.tsr.SelectTaskStagesByPcode(ctx, code)
}

func (repo *TaskStagesRepository) CreateTaskStageByStrcut(ctx context.Context, stage *model.TaskStages) error {
	return repo.tsr.InsertTaskStages(ctx, stage)
}

func (repo *TaskStagesRepository) DeleteTaskStagesByCode(ctx context.Context, code int64) error {
	return repo.tsr.DeleteTaskStagesById(ctx, code)
}

func (repo *TaskStagesRepository) EditTaskStagesByStruct(ctx context.Context, code int64, ts *model.TaskStages) error {
	return repo.tsr.UpdateTaskStagesByStruct(ctx, code, ts)
}
