package dao

import (
	"context"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository/database"
)

type TaskWorkTimeDao struct {
	db *database.GormConn
}

func (dao *TaskWorkTimeDao) InsertTaskWorkTime(ctx context.Context, twt *model.TaskWorkTime) error {
	session := dao.db.Session(ctx)
	err := session.Create(&twt).Error
	return err
}

func (dao *TaskWorkTimeDao) SelectTaskWorkTimeList(ctx context.Context, taskCode int64) (list []*model.TaskWorkTime, err error) {
	session := dao.db.Session(ctx)
	err = session.Model(&model.TaskWorkTime{}).Where("task_code = ?", taskCode).Find(&list).Error
	return
}

func NewTaskWorkTimeDao() *TaskWorkTimeDao {
	return &TaskWorkTimeDao{
		db: database.NewGormSession(),
	}
}
