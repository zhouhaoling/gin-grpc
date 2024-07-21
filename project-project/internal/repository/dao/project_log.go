package dao

import (
	"context"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository/database"
)

type ProjectLogDao struct {
	db *database.GormConn
}

func (dao *ProjectLogDao) SelectProjectLogByMid(ctx context.Context, mid int64, page int64, size int64) (list []*model.ProjectLog, total int64, err error) {
	session := dao.db.Session(ctx)
	offset := (page - 1) * size
	err = session.Model(&model.ProjectLog{}).
		Where("member_code = ?", mid).
		Limit(int(size)).
		Offset(int(offset)).Order("create_time desc").Find(&list).Error
	err = session.Model(&model.ProjectLog{}).
		Where("member_code = ?", mid).Count(&total).Error
	return
}

func (dao *ProjectLogDao) SelectProjectLogByTaskCode(ctx context.Context, taskCode int64, comment int) (list []*model.ProjectLog, total int64, err error) {
	session := dao.db.Session(ctx)
	pl := session.Model(&model.ProjectLog{})
	if comment == 1 {
		err = pl.Where("source_code=? and is_comment=?", taskCode, comment).Find(&list).Error
		pl.Where("source_code=? and is_comment=?", taskCode, comment).Count(&total)
	} else {
		err = pl.Where("source_code=?", taskCode).Find(&list).Error
		pl.Where("source_code=?", taskCode).Count(&total)
	}
	return
}

func (dao *ProjectLogDao) SelectProjectLogByTaskCodePage(ctx context.Context, taskCode int64, comment int, page int, pageSize int) (list []*model.ProjectLog, total int64, err error) {
	session := dao.db.Session(ctx)
	pl := session.Model(&model.ProjectLog{})
	offset := (page - 1) * pageSize
	if comment == 1 {
		err = pl.Where("source_code=? and is_comment=?", taskCode, comment).Limit(pageSize).Offset(offset).Find(&list).Error
		pl.Where("source_code=? and is_comment=?", taskCode, comment).Count(&total)
	} else {
		err = pl.Where("source_code=?", taskCode).Limit(pageSize).Offset(offset).Find(&list).Error
		pl.Where("source_code=?", taskCode).Count(&total)
	}
	return
}

func (dao *ProjectLogDao) InsertProjectLogByStruct(ctx context.Context, pl *model.ProjectLog) error {
	session := dao.db.Session(ctx)
	err := session.Create(&pl).Error
	return err
}

func NewProjectLogDao() *ProjectLogDao {
	return &ProjectLogDao{
		db: database.NewGormSession(),
	}
}
