package repository

import "test.com/project-user/internal/repository/dao/mysql"

type DingTalkRepository struct {
	dt *mysql.DingTalkDao
}

func NewDingTalkRepository() *DingTalkRepository {
	return &DingTalkRepository{
		dt: mysql.NewDingTalkDao(),
	}
}
