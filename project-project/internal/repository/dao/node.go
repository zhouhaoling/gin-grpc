package dao

import (
	"context"

	"test.com/project-project/internal/model"

	"test.com/project-project/internal/repository/database"
)

type ProjectNodeDao struct {
	conn *database.GormConn
}

func (m *ProjectNodeDao) FindAll(ctx context.Context) (pms []*model.ProjectNode, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&model.ProjectNode{}).Find(&pms).Error
	return
}

func NewProjectNodeDao() *ProjectNodeDao {
	return &ProjectNodeDao{
		conn: database.NewGormSession(),
	}
}
