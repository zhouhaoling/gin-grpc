package repo

import (
	"context"

	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/database"
)

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	InsertMember(conn database.DBConn, ctx context.Context, member *model.Member) error
	SelectMemberByAccount(ctx context.Context, account string) (model.Member, error)
	SelectMemberByMemId(ctx context.Context, mid int64) (*model.Member, error)
	SelectMemberByMemIdAndName(ctx context.Context, mid int64, name string) (bool, error)
	SelectMemberByMemIds(ctx context.Context, ids []int64) (list []*model.Member, err error)
}
