package dao

import (
	"context"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository/database"
)

type TaskStagesTemplateDao struct {
	conn *database.GormConn
}

func (dao *TaskStagesTemplateDao) InsertTaskStagesTemplate(ctx context.Context, dbConn database.DBConn, list []*model.MsTaskStagesTemplate) error {
	dao.conn = dbConn.(*database.GormConn)
	err := dao.conn.Tx(ctx).Where(&model.MsTaskStagesTemplate{}).Create(&list).Error
	return err
}

func (dao *TaskStagesTemplateDao) SelectTaskStagesTemplateByPtcodes(ctx context.Context, ptcodes []int64) (list []*model.MsTaskStagesTemplate, err error) {
	session := dao.conn.Session(ctx)
	err = session.Model(&model.MsTaskStagesTemplate{}).
		Where("project_template_code in (?)", ptcodes).
		Find(&list).Error
	return list, err
}

// SelectTaskTemplateByProjectTemplateId 根据template_project_code查询任务模板, tpId = template_project_code
func (dao *TaskStagesTemplateDao) SelectTaskTemplateByProjectTemplateId(ctx context.Context, tpId int) (list []*model.MsTaskStagesTemplate, err error) {
	session := dao.conn.Session(ctx)
	err = session.Model(&model.MsTaskStagesTemplate{}).
		Where("project_template_code = ?", tpId).
		Order("sort desc, id asc").
		Find(&list).Error
	return list, err
}

func (dao *TaskStagesTemplateDao) SelectInProTemIds(ctx context.Context, ids []int) (mtst []model.MsTaskStagesTemplate, err error) {
	session := dao.conn.Session(ctx)
	err = session.Where("project_template_code in ?", ids).Find(&mtst).Error
	return mtst, err
}

func NewTaskStagesTemplateDao() *TaskStagesTemplateDao {
	return &TaskStagesTemplateDao{
		conn: database.NewGormSession(),
	}
}
