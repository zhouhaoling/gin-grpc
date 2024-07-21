package repository

import (
	"test.com/project-user/internal/repo"
	"test.com/project-user/internal/repository/dao"
)

type DingTalkRepository struct {
	dt repo.DingTalkRepo
}

func NewDingTalkRepository() *DingTalkRepository {
	return &DingTalkRepository{
		dt: dao.NewDingTalkDao(),
	}
}
