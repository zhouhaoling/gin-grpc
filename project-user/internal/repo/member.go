package repo

import (
	"context"

	"test.com/project-user/internal/repository/database"

	"test.com/project-user/internal/model"
)

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	InsertMember(conn database.DBConn, ctx context.Context, member *model.Member) error
}
