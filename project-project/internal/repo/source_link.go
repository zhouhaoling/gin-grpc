package repo

import (
	"context"

	"test.com/project-project/internal/repository/database"

	"test.com/project-project/internal/model"
)

type SourceLinkRepo interface {
	InsertSourceLinkByStruct(ctx context.Context, conn database.DBConn, sl *model.SourceLink) error
	SelectSourceLinkByTaskCode(ctx context.Context, taskCode int64) (list []*model.SourceLink, err error)
}
