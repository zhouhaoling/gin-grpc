package repository

import (
	"context"
	"time"

	"test.com/common/errs"
	"test.com/project-project/config"

	"test.com/project-project/internal/repository/database"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository/dao"
)

type ProjectRepository struct {
	pjt repo.ProjectRepo
	plt repo.ProjectLogRepo
}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{
		pjt: dao.NewProjectDao(),
		plt: dao.NewProjectLogDao(),
	}
}

var (
	AndArchive = "and archive = 1"
	AndDeleted = "and deleted = 1"
	AndMy      = "and deleted = 0"
)

func (repo *ProjectRepository) FindProjectByMemId(ctx context.Context, mid, page, size int64, flag int) (pms []*model.ProjectAndMember, total int64, err error) {
	switch flag {
	case 0:
		pms, total, err = repo.pjt.SelectProjectByMemId(ctx, mid, page, size, AndMy)
	case 1:
		pms, total, err = repo.pjt.SelectProjectByMemId(ctx, mid, page, size, AndArchive)
	case 2:
		pms, total, err = repo.pjt.SelectProjectByMemId(ctx, mid, page, size, AndDeleted)
	case 3:
		pms, total, err = repo.pjt.SelectCollectProjectByMemId(ctx, mid, page, size)

	}
	return
}

func (repo *ProjectRepository) CreateProject(conn database.DBConn, ctx context.Context, project *model.Project) error {
	return repo.pjt.InsertProject(conn, ctx, project)
}

func (repo *ProjectRepository) CreateProjectMember(conn database.DBConn, ctx context.Context, pm *model.ProjectMember) error {
	return repo.pjt.InsertProjectMember(conn, ctx, pm)
}

func (repo *ProjectRepository) FindProjectByPidAndMemberId(ctx context.Context, pid, mid int64) (*model.ProjectAndMember, error) {
	return repo.pjt.SelectProjectByPidAndMemberId(ctx, pid, mid)
}

func (repo *ProjectRepository) FindCollectByPidAndMemId(ctx context.Context, pid, mid int64) (bool, error) {
	return repo.pjt.SelectCollectByPidAndMemId(ctx, pid, mid)
}

func (repo *ProjectRepository) RecycleProjectByPid(ctx context.Context, pid int64, deleted bool) error {
	return repo.pjt.UpdateProjectByPid(ctx, pid, deleted)
}

func (repo *ProjectRepository) CollectProjectByType(ctx context.Context, pid, mid int64, collect string) error {
	var err error
	switch collect {
	case config.Collect:
		pc := &model.Collection{
			ProjectCode: pid,
			MemberCode:  mid,
			CreateTime:  time.Now().UnixMilli(),
		}
		err = repo.pjt.InsertProjectCollect(ctx, pc)
	case config.DCollect:
		err = repo.pjt.DeleteProjectCollect(ctx, pid, mid)
	default:
		err = errs.GrpcError(model.ErrorServerBusy)
	}

	return err
}

func (repo *ProjectRepository) UpdateProjectByStruct(ctx context.Context, pid int64, project *model.Project) error {
	return repo.pjt.UpdatesProjectByStructAndPid(ctx, pid, project)
}

func (repo *ProjectRepository) FindProjectMemberByPid(ctx context.Context, pid int64) ([]*model.ProjectMember, int64, error) {
	return repo.pjt.SelectProjectMemberByPid(ctx, pid)
}

func (repo *ProjectRepository) FindProjectByPId(ctx context.Context, pid int64) (*model.Project, error) {
	return repo.pjt.SelectProjectByPid(ctx, pid)
}

func (repo *ProjectRepository) FindProjectByPIds(ctx context.Context, pids []int64) ([]*model.Project, error) {
	return repo.pjt.SelectProjectByPids(ctx, pids)
}

func (repo *ProjectRepository) SaveProjectLog(ctx context.Context, pl *model.ProjectLog) error {
	return repo.plt.InsertProjectLogByStruct(ctx, pl)
}

func (repo *ProjectRepository) FindProjectLogByTaskCode(ctx context.Context, taksCode int64, comment int) ([]*model.ProjectLog, int64, error) {
	return repo.plt.SelectProjectLogByTaskCode(ctx, taksCode, comment)
}

func (repo *ProjectRepository) FindProjectLogByTaskCodePage(ctx context.Context, taskCode int64, comment int, page int, pageSize int) ([]*model.ProjectLog, int64, error) {
	return repo.plt.SelectProjectLogByTaskCodePage(ctx, taskCode, comment, page, pageSize)
}

func (repo *ProjectRepository) FindProjectLogByMid(ctx context.Context, mid int64, page int64, pageSize int64) ([]*model.ProjectLog, int64, error) {
	return repo.plt.SelectProjectLogByMid(ctx, mid, page, pageSize)
}
