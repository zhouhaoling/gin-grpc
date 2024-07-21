package repository

import (
	"context"

	"test.com/project-project/internal/repository/database"

	"test.com/project-project/internal/model"

	pg "test.com/project-grpc/project_grpc"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository/dao"
)

type ProjectTemplateRepository struct {
	prt  repo.ProjectTemplateRepo
	tstr repo.TaskStagesTemplateRepo
	tsr  repo.TaskStagesRepo
}

func NewProjectTemplateRepository() *ProjectTemplateRepository {
	return &ProjectTemplateRepository{
		prt:  dao.NewProjectTemplateDao(),
		tstr: dao.NewTaskStagesTemplateDao(),
		tsr:  dao.NewTaskStagesDao(),
	}
}

func (repo *ProjectTemplateRepository) FindProjectTemplate(ctx context.Context, msg *pg.ProjectRequest, orgId int64) (pts []model.ProjectTemplate, total int64, err error) {
	switch msg.ViewType {
	case -1:
		pts, total, err = repo.prt.SelectProjectTemplateAll(ctx, orgId, msg.Page, msg.PageSize)
	case 0:
		pts, total, err = repo.prt.SelectProjectTemplateCustom(ctx, msg.MemberId, orgId, msg.Page, msg.PageSize)
	case 1:
		pts, total, err = repo.prt.SelectProjectTemplateSystem(ctx, msg.Page, msg.PageSize)
	}
	return
}

func (repo *ProjectTemplateRepository) FindInProTemIds(ctx context.Context, ids []int) ([]model.MsTaskStagesTemplate, error) {
	return repo.tstr.SelectInProTemIds(ctx, ids)
}

func (repo *ProjectTemplateRepository) FindTaskTemplateByProjectTemplateId(ctx context.Context, tpId int) ([]*model.MsTaskStagesTemplate, error) {
	return repo.tstr.SelectTaskTemplateByProjectTemplateId(ctx, tpId)
}

func (repo *ProjectTemplateRepository) FindProjectTemplateByMid(ctx context.Context, mid int64) ([]*model.ProjectTemplate, error) {
	return repo.prt.SelectProjectTemplateByMid(ctx, mid)
}

func (repo *ProjectTemplateRepository) CreateProjectTemplate(ctx context.Context, conn database.DBConn, list []*model.ProjectTemplate) error {
	return repo.prt.InsertProjectTemplate(ctx, conn, list)
}

func (repo *ProjectTemplateRepository) FindTaskStagesTemplateByPtcodes(ctx context.Context, ptcodes []int64) ([]*model.MsTaskStagesTemplate, error) {
	return repo.tstr.SelectTaskStagesTemplateByPtcodes(ctx, ptcodes)
}

func (repo *ProjectTemplateRepository) CreateTaskStagesTemplate(ctx context.Context, conn database.DBConn, list []*model.MsTaskStagesTemplate) error {
	return repo.tstr.InsertTaskStagesTemplate(ctx, conn, list)
}
