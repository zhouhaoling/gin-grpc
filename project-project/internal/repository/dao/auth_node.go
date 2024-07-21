package dao

import (
	"context"

	"test.com/project-project/internal/model"

	"test.com/project-project/internal/repository/database"
)

type ProjectAuthNodeDao struct {
	db *database.GormConn
}

func (p *ProjectAuthNodeDao) InsertProjectAuthNodeList(ctx context.Context, conn database.DBConn, list []*model.ProjectAuthNode) error {
	p.db = conn.(*database.GormConn)
	err := p.db.Tx(ctx).Create(list).Error
	return err
}

func (p *ProjectAuthNodeDao) SelectProjectAuthNode(ctx context.Context) (list []*model.ProjectAuthNode, err error) {
	session := p.db.Session(ctx)
	err = session.Model(&model.ProjectAuthNode{}).Where("auth = 2").Find(&list).Error
	return
}

func (p *ProjectAuthNodeDao) DeleteByAuthId(ctx context.Context, conn database.DBConn, authId int64) error {
	p.db = conn.(*database.GormConn)
	tx := p.db.Tx(ctx)
	err := tx.Where("auth = ?", authId).Delete(&model.ProjectAuthNode{}).Error
	return err
}

func (p *ProjectAuthNodeDao) Save(ctx context.Context, conn database.DBConn, authId int64, nodes []string) error {
	p.db = conn.(*database.GormConn)
	tx := p.db.Tx(ctx)
	var list []*model.ProjectAuthNode
	for _, v := range nodes {
		pn := &model.ProjectAuthNode{
			Auth: authId,
			Node: v,
		}
		list = append(list, pn)
	}
	err := tx.Create(list).Error
	return err
}

func (p *ProjectAuthNodeDao) FindNodeStringList(ctx context.Context, authId int64) (list []string, err error) {
	session := p.db.Session(ctx)
	err = session.Model(&model.ProjectAuthNode{}).Where("auth=?", authId).Select("node").Find(&list).Error
	return
}

func NewProjectAuthNodeDao() *ProjectAuthNodeDao {
	return &ProjectAuthNodeDao{
		db: database.NewGormSession(),
	}
}
