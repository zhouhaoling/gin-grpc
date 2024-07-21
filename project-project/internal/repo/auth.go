package repo

import (
	"context"

	"test.com/project-project/internal/model"

	"test.com/project-project/internal/repository/database"
)

type ProjectAuthNodeRepo interface {
	FindNodeStringList(ctx context.Context, authId int64) ([]string, error)
	DeleteByAuthId(ctx context.Context, conn database.DBConn, id int64) error
	Save(ctx context.Context, conn database.DBConn, id int64, nodes []string) error
	SelectProjectAuthNode(ctx context.Context) (list []*model.ProjectAuthNode, err error)
	InsertProjectAuthNodeList(ctx context.Context, conn database.DBConn, list []*model.ProjectAuthNode) error
}
