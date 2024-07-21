package dao

import (
	"context"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository/database"
)

type TaskStagesDao struct {
	db *database.GormConn
}

func (dao *TaskStagesDao) UpdateTaskStagesByStruct(ctx context.Context, code int64, ts *model.TaskStages) error {
	session := dao.db.Session(ctx)
	err := session.Model(&model.TaskStages{}).Where("id = ?", code).Updates(ts).Error
	return err
}

func (dao *TaskStagesDao) DeleteTaskStagesById(ctx context.Context, code int64) error {
	session := dao.db.Session(ctx)
	err := session.Delete(&model.TaskStages{}, code).Error
	return err
}

func (dao *TaskStagesDao) InsertTaskStages(ctx context.Context, stage *model.TaskStages) error {
	session := dao.db.Session(ctx)
	err := session.Create(stage).Error
	return err
}

func (dao *TaskStagesDao) SelectTaskStagesByPcode(ctx context.Context, code int64) (ts *model.TaskStages, err error) {
	session := dao.db.Session(ctx)
	err = session.Where("project_code = ?", code).Order("sort desc").First(&ts).Error
	return
}

func (dao *TaskStagesDao) SelectByStageId(ctx context.Context, sid int) (ts *model.TaskStages, err error) {
	err = dao.db.Session(ctx).Where("id = ?", sid).First(&ts).Error
	return
}

func (dao *TaskStagesDao) SelectTaskStagesByProjectId(ctx context.Context, pid, page, pageSize int64) (list []*model.TaskStages, total int64, err error) {
	session := dao.db.Session(ctx)
	err = session.Where("project_code = ?", pid).
		Order("sort asc").
		Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).
		Find(&list).Error
	session.Model(&model.TaskStages{}).Where("project_code = ?", pid).Count(&total)
	return
}

func (dao *TaskStagesDao) SaveTaskStages(conn database.DBConn, ctx context.Context, ts *model.TaskStages) error {
	dao.db = conn.(*database.GormConn)
	err := dao.db.Tx(ctx).Save(ts).Error
	return err
}

func NewTaskStagesDao() *TaskStagesDao {
	return &TaskStagesDao{
		db: database.NewGormSession(),
	}
}
