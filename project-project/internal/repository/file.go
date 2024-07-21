package repository

import (
	"context"

	"test.com/project-project/internal/repository/database"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository/dao"
)

type FileRepository struct {
	fr repo.FileRepo
	sl repo.SourceLinkRepo
}

func (repo *FileRepository) SaveFileInfo(ctx context.Context, conn database.DBConn, f *model.File) error {
	return repo.fr.InsertFileByStruct(ctx, conn, f)
}

func (repo *FileRepository) SaveSourceLinkInfo(ctx context.Context, conn database.DBConn, sl *model.SourceLink) error {
	return repo.sl.InsertSourceLinkByStruct(ctx, conn, sl)
}

func (repo *FileRepository) FindSourceLinksByTaskCode(ctx context.Context, code int64) ([]*model.SourceLink, error) {
	return repo.sl.SelectSourceLinkByTaskCode(ctx, code)
}

func (repo *FileRepository) FindFileByIds(ctx context.Context, list []int64) ([]*model.File, error) {
	return repo.fr.FindFileByIds(ctx, list)
}

func NewFileRepository() *FileRepository {
	return &FileRepository{
		fr: dao.NewFileDao(),
		sl: dao.NewSourceLinkDao(),
	}
}
