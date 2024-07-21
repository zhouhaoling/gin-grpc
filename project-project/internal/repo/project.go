package repo

import (
	"context"

	"test.com/project-project/internal/repository/database"

	"test.com/project-project/internal/model"
)

type ProjectRepo interface {
	SelectProjectByMemId(ctx context.Context, memId int64, page int64, size int64, condition string) ([]*model.ProjectAndMember, int64, error)
	SelectCollectProjectByMemId(ctx context.Context, mid int64, page int64, size int64) ([]*model.ProjectAndMember, int64, error)
	InsertProject(conn database.DBConn, ctx context.Context, project *model.Project) error
	InsertProjectMember(conn database.DBConn, ctx context.Context, pm *model.ProjectMember) error
	SelectProjectByPidAndMemberId(ctx context.Context, pid int64, mid int64) (*model.ProjectAndMember, error)
	SelectCollectByPidAndMemId(ctx context.Context, pid int64, mid int64) (bool, error)
	UpdateProjectByPid(ctx context.Context, pid int64, deleted bool) error
	InsertProjectCollect(ctx context.Context, pc *model.Collection) error
	DeleteProjectCollect(ctx context.Context, pid int64, mid int64) error
	UpdatesProjectByStructAndPid(ctx context.Context, pid int64, project *model.Project) error
	SelectProjectMemberByPid(ctx context.Context, pid int64) (list []*model.ProjectMember, total int64, err error)
	SelectProjectByPid(ctx context.Context, pid int64) (pj *model.Project, err error)
	SelectProjectByPids(ctx context.Context, pids []int64) (pList []*model.Project, err error)
}
