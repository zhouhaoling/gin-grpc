package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository/database"
)

type AccountDao struct {
	db *database.GormConn
}

func (dao *AccountDao) InsertMemberAccountByStruct(conn database.DBConn, ctx context.Context, ma *model.MemberAccount) error {
	dao.db = conn.(*database.GormConn)
	err := dao.db.Tx(ctx).Create(&ma).Error
	return err
}

func (dao *AccountDao) SelectAccountByMemberId(ctx context.Context, mid int64) (account *model.MemberAccount, err error) {
	session := dao.db.Session(ctx)
	err = session.Where("member_code=?", mid).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func (dao *AccountDao) FindList(ctx context.Context, condition string, organizationCode int64, departmentCode int64, page int64, pageSize int64) (list []*model.MemberAccount, total int64, err error) {
	session := dao.db.Session(ctx)
	offset := (page - 1) * pageSize
	err = session.Model(&model.MemberAccount{}).
		Where("organization_code=?", organizationCode).
		Where(condition).Limit(int(pageSize)).Offset(int(offset)).Find(&list).Error
	err = session.Model(&model.MemberAccount{}).
		Where("organization_code=?", organizationCode).
		Where(condition).Count(&total).Error
	return
}

func NewAccountDao() *AccountDao {
	return &AccountDao{db: database.NewGormSession()}
}
