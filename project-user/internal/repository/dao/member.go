package dao

import (
	"context"

	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/database"
)

type MemberDao struct {
	db *database.GormConn
}

func (dao *MemberDao) SelectMemberByMemIds(ctx context.Context, ids []int64) (list []*model.Member, err error) {
	if len(ids) <= 0 {
		return nil, nil
	}
	session := dao.db.Session(ctx)
	err = session.Model(&model.Member{}).Where("mid in (?)", ids).Find(&list).Error
	return list, err
}

func (dao *MemberDao) SelectMemberByMemIdAndName(ctx context.Context, mid int64, name string) (bool, error) {
	var count int64
	err := dao.db.Session(ctx).Model(&model.Member{}).Where("mid = ? and name = ? and status = 1", mid, name).Count(&count).Error
	return count > 0, err
}

func (dao *MemberDao) SelectMemberByMemId(ctx context.Context, mid int64) (member *model.Member, err error) {
	session := dao.db.Session(ctx)
	err = session.Where("mid = ?", mid).First(&member).Error
	return member, err
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

// GetMemberByMobile 判断手机号是否存在
func (dao *MemberDao) GetMemberByMobile(ctx context.Context, mobile string) (bool, error) {
	//TODO implement me
	var count int64

	err := dao.db.Session(ctx).Model(&model.Member{}).Where("mobile = ?", mobile).Count(&count).Error
	return count > 0, err
}

func (dao *MemberDao) SelectMemberByAccount(ctx context.Context, account string) (member model.Member, err error) {
	result := dao.db.Session(ctx).Where("account = ?", account).First(&member)
	return member, result.Error
}
