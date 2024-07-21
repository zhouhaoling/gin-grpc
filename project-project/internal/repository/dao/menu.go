package dao

import (
	"context"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository/database"
)

type MenuDao struct {
	conn *database.GormConn
}

func (dao *MenuDao) SelectMenus(ctx context.Context) (pms []*model.ProjectMenu, err error) {
	session := dao.conn.Session(ctx)
	err = session.Order("pid,sort asc, id asc").Find(&pms).Error
	return
}

func NewMenuDao() *MenuDao {
	return &MenuDao{
		conn: database.NewGormSession(),
	}
}
