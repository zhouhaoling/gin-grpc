package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"test.com/project-project/config"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository/database"
)

type TaskDao struct {
	db *database.GormConn
}

func (dao *TaskDao) SelectTaskByIds(ctx context.Context, taskIds []int64) (list []*model.Task, err error) {
	session := dao.db.Session(ctx)
	err = session.Model(&model.Task{}).Where("id in (?)", taskIds).Find(&list).Error
	return
}

func (dao *TaskDao) SelectTaskMemberByTaskIdAndPage(ctx context.Context, code int64, page int64, pageSize int64) (list []*model.TaskMember, total int64, err error) {
	session := dao.db.Session(ctx)
	offset := (page - 1) * pageSize
	err = session.Model(&model.TaskMember{}).Where("task_code = ?", code).Limit(int(pageSize)).Offset(int(offset)).Find(&list).Error
	session.Model(&model.TaskMember{}).Where("task_code = ?", code).Count(&total)
	return
}

func (dao *TaskDao) SelectTaskByAssignTo(ctx context.Context, mid int64, done int, page int64, pageSize int64) (list []*model.Task, total int64, err error) {
	session := dao.db.Session(ctx)
	offset := (page - 1) * pageSize
	err = session.Model(&model.Task{}).Where("assign_to=? and deleted=0 and done=?", mid, done).Limit(int(pageSize)).Offset(int(offset)).Find(&list).Error
	err = session.Model(&model.Task{}).Where("assign_to=? and deleted=0 and done=?", mid, done).Count(&total).Error
	return
}

func (dao *TaskDao) SelectTaskMemberByMemberCode(ctx context.Context, mid int64, done int, page int64, pageSize int64) (list []*model.Task, total int64, err error) {
	session := dao.db.Session(ctx)
	offset := (page - 1) * pageSize
	sql := "select a.* from ms_task a,ms_task_member b where a.id=b.task_code and member_code=? and a.deleted=0 and a.done=? limit ?, ?"
	raw := session.Model(&model.Task{}).Raw(sql, mid, done, offset, pageSize)
	err = raw.Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}
	sqlCount := "select count(*) from ms_task a,ms_task_member b where a.id=b.task_code and member_code=? and a.deleted=0 and a.done=?"
	rawCount := session.Model(&model.Task{}).Raw(sqlCount, mid, done)
	err = rawCount.Scan(&total).Error
	return
}

func (dao *TaskDao) SelectTaskByCreateBy(ctx context.Context, mid int64, done int, page int64, pageSize int64) (list []*model.Task, total int64, err error) {
	session := dao.db.Session(ctx)
	offset := (page - 1) * pageSize
	err = session.Model(&model.Task{}).Where("create_by=? and deleted=0 and done=?", mid, done).Limit(int(pageSize)).Offset(int(offset)).Find(&list).Error
	err = session.Model(&model.Task{}).Where("create_by=? and deleted=0 and done=?", mid, done).Count(&total).Error
	return
}

func (dao *TaskDao) SelectTaskByStageCodeLtSort(ctx context.Context, stageCode int, sort int) (ts *model.Task, err error) {
	session := dao.db.Session(ctx)
	err = session.Where("stage_code = ? and sort < ?", stageCode, sort).Order("sort desc").Limit(1).First(&ts).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (dao *TaskDao) UpdateTaskSort(ctx context.Context, conn database.DBConn, task *model.Task) error {
	dao.db = conn.(*database.GormConn)
	err := dao.db.Tx(ctx).Where("id = ?", task.Id).Select("sort", "stage_code").Updates(&task).Error
	return err
}

func (dao *TaskDao) SelectTaskById(ctx context.Context, taskId int64) (task *model.Task, err error) {
	session := dao.db.Session(ctx)
	err = session.Where("id = ?", taskId).First(&task).Error
	return
}

func (dao *TaskDao) InsertTask(conn database.DBConn, ctx context.Context, ts *model.Task) error {
	dao.db = conn.(*database.GormConn)
	return dao.db.Tx(ctx).Save(&ts).Error
}

func (dao *TaskDao) InsertTaskMember(conn database.DBConn, ctx context.Context, tm *model.TaskMember) error {
	dao.db = conn.(*database.GormConn)
	return dao.db.Tx(ctx).Save(&tm).Error
}

func (dao *TaskDao) SelectTaskMemberByTaskId(ctx context.Context, taskId int64, mid int64) (task *model.TaskMember, err error) {
	session := dao.db.Session(ctx)
	result := session.Where("stage_code = ? and member_code = ?", taskId, mid).
		Limit(1).
		Find(&task)
	err = result.Error
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return
}

func (dao *TaskDao) SelectTaskSortByPIdAndSId(ctx context.Context, pid int64, sid int64) (v *int, err error) {
	session := dao.db.Session(ctx)
	err = session.Model(&model.Task{}).Where("project_code = ? and stage_code=?", pid, sid).
		Select("max(sort) as sort").
		Scan(&v).Error
	return
}

func (dao *TaskDao) SelectTaskMaxIdNum(ctx context.Context, pid int64) (v *int, err error) {
	session := dao.db.Session(ctx)
	err = session.Model(&model.Task{}).Where("project_code = ?", pid).
		Select("max(id_num) as maxIdNum").
		Scan(&v).Error
	return
}

func (dao *TaskDao) SelectTaskByStageCode(ctx context.Context, stageCode int) (list []*model.Task, err error) {
	//select * from ms_task where stage_code = 1 and delete = 0 order by sort asc
	session := dao.db.Session(ctx)
	err = session.Model(&model.Task{}).Where("stage_code = ? and deleted = ?", stageCode, config.NoDeleted).
		Order("sort asc").
		Find(&list).Error
	return
}

func NewTaskDao() *TaskDao {
	return &TaskDao{
		db: database.NewGormSession(),
	}
}
