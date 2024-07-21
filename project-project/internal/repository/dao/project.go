package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"test.com/project-project/config"

	"test.com/project-project/internal/repository/database"

	"test.com/project-project/internal/model"
)

type ProjectDao struct {
	db *database.GormConn
}

func (dao *ProjectDao) SelectProjectByPids(ctx context.Context, pids []int64) (pList []*model.Project, err error) {
	session := dao.db.Session(ctx)
	err = session.Model(&model.Project{}).Where("id in (?)", pids).Find(&pList).Error
	return
}

func (dao *ProjectDao) SelectProjectByPid(ctx context.Context, pid int64) (pj *model.Project, err error) {
	err = dao.db.Session(ctx).Where("id = ?", pid).First(&pj).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (dao *ProjectDao) SelectProjectMemberByPid(ctx context.Context, pid int64) (list []*model.ProjectMember, total int64, err error) {
	session := dao.db.Session(ctx)
	err = session.Where("project_code = ?", pid).Find(&list).Error
	session.Model(&model.ProjectMember{}).Where("project_code = ?", pid).Count(&total)
	return
}

func (dao *ProjectDao) UpdatesProjectByStructAndPid(ctx context.Context, pid int64, project *model.Project) error {
	session := dao.db.Session(ctx)
	err := session.Model(&model.Project{}).Where("id = ?", pid).Updates(project).Error
	return err
}

func (dao *ProjectDao) InsertProjectCollect(ctx context.Context, pc *model.Collection) error {
	session := dao.db.Session(ctx)
	return session.Create(pc).Error
}

func (dao *ProjectDao) DeleteProjectCollect(ctx context.Context, pid int64, mid int64) error {
	return dao.db.Session(ctx).Where("member_code = ? and project_code = ?", mid, pid).Delete(&model.Collection{}).Error
}

func (dao *ProjectDao) UpdateProjectByPid(ctx context.Context, pid int64, deleted bool) error {
	session := dao.db.Session(ctx)
	var project model.Project
	deleteTime := time.Now().UnixMilli()
	var err error
	if deleted {
		err = session.Model(&project).Where("id = ?", pid).Updates(model.Project{Deleted: config.Deleted, DeletedTime: deleteTime}).Error
	} else {
		err = session.Model(&project).Where("id = ?", pid).Updates(map[string]interface{}{"deleted": config.NoDeleted, "deleted_time": sql.NullInt64{}}).Error
	}
	return err
}

func (dao *ProjectDao) SelectCollectByPidAndMemId(ctx context.Context, pid int64, mid int64) (bool, error) {
	var count int64
	session := dao.db.Session(ctx)
	sql := fmt.Sprintf("select count(*) from ms_project_collection where project_code = ? and member_code = ? ")
	raw := session.Raw(sql, pid, mid)
	err := raw.Scan(&count).Error
	return count > 0, err
}

func (dao *ProjectDao) SelectProjectByPidAndMemberId(ctx context.Context, pid int64, mid int64) (pam *model.ProjectAndMember, err error) {
	session := dao.db.Session(ctx)
	sql := fmt.Sprintf("select a.*, b.project_code, b.member_code, b.join_time, b.is_owner, b.authorize from ms_project a, ms_project_member b where a.id = b.project_code and b.member_code = ? and b.project_code = ? limit 1")
	raw := session.Raw(sql, mid, pid)
	err = raw.Scan(&pam).Error
	return pam, err
}

func (dao *ProjectDao) InsertProjectMember(conn database.DBConn, ctx context.Context, pm *model.ProjectMember) error {
	dao.db = conn.(*database.GormConn)
	return dao.db.Tx(ctx).Create(pm).Error
}

func (dao *ProjectDao) InsertProject(conn database.DBConn, ctx context.Context, project *model.Project) error {
	dao.db = conn.(*database.GormConn)
	return dao.db.Tx(ctx).Create(project).Error
}

func (dao *ProjectDao) SelectCollectProjectByMemId(ctx context.Context, mid int64, page int64, size int64) (pm []*model.ProjectAndMember, total int64, err error) {
	//TODO implement me
	session := dao.db.Session(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from ms_project where id in (select project_code from ms_project_collection where member_code = ?) order by sort limit ?, ?")
	db := session.Raw(sql, mid, index, size)
	db.Scan(&pm)
	query := fmt.Sprintf("member_code = ?")
	err = session.Model(&model.Collection{}).Where(query, mid).Count(&total).Error
	return
}

func (dao *ProjectDao) SelectProjectByMemId(ctx context.Context, memId, page, size int64, condition string) (pm []*model.ProjectAndMember, total int64, err error) {
	session := dao.db.Session(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from ms_project a, ms_project_member b where a.id = b.project_code and b.member_code = ? %s order by sort limit ?,?", condition)
	db := session.Raw(sql, memId, index, size)
	db.Scan(&pm)
	query := fmt.Sprintf("select count(*) from ms_project a, ms_project_member b where a.id = b.project_code and b.member_code = ? %s", condition)
	tx := session.Raw(query, memId)
	err = tx.Scan(&total).Error
	return pm, total, err
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{
		db: database.NewGormSession(),
	}
}
