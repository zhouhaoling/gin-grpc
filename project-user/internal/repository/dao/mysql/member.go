package mysql

import "context"

type MemberDao struct {
}

func (m MemberDao) GetMemberByEmail(ctx context.Context, email string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m MemberDao) GetMemberByMobile(ctx context.Context, mobile string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func NewMemberDao() *MemberDao {
	return &MemberDao{}
}
