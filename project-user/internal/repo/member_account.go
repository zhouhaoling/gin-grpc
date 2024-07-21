package repo

import (
	"context"

	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/database"
)

type MemberAccountRepo interface {
	InsertMemberAccountByStruct(conn database.DBConn, ctx context.Context, ma *model.MemberAccount) error
	UpdateMemberAccountByLastTime(ctx context.Context, mid int64) error
}
