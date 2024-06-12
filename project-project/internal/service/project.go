package service

import (
	"context"

	"test.com/project-project/internal/repository"
	"test.com/project-project/internal/repository/dao"

	"test.com/project-project/internal/repository/database"

	pg "test.com/project-grpc/project_grpc"
	"test.com/project-project/internal/repo"
)

type ProjectService struct {
	pg.UnimplementedProjectServiceServer
	cache repo.Cache
	tran  database.Transaction
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		cache: dao.RC,
		tran:  repository.NewTransaction(),
	}
}

func (p ProjectService) Index(ctx context.Context, msg *pg.IndexRequest) (*pg.IndexResponse, error) {
	//TODO implement me
	return nil, nil
}
