package mysql

import (
	"context"

	"test.com/project-user/internal/model"

	"test.com/project-user/internal/repository/database"
)

type MemberDao struct {
	db *database.GormConn
}

func NewMemberDao() *MemberDao {
	return &MemberDao{
		db: database.NewGormSession(),
	}
}

func (dao *MemberDao) InsertMember(conn database.DBConn, ctx context.Context, member *model.Member) error {
	//return dao.db.Session(ctx).Create(member).Error
	dao.db = conn.(*database.GormConn)
	return dao.db.Tx(ctx).Create(member).Error
}

func (dao *MemberDao) GetMemberByEmail(ctx context.Context, email string) (bool, error) {
	var count int64

	err := dao.db.Session(ctx).Model(&model.Member{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (dao *MemberDao) GetMemberByMobile(ctx context.Context, mobile string) (bool, error) {
	//TODO implement me
	var count int64

	err := dao.db.Session(ctx).Model(&model.Member{}).Where("mobile = ?", mobile).Count(&count).Error
	return count > 0, err
}
