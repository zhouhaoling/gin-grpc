package repo

import (
	"context"

	"test.com/project-project/internal/repository/database"

	"test.com/project-project/internal/model"
)

type AccountRepo interface {
	FindList(ctx context.Context, condition string, organizationCode int64, departmentCode int64, page int64, pageSize int64) ([]*model.MemberAccount, int64, error)
	SelectAccountByMemberId(ctx context.Context, mid int64) (account *model.MemberAccount, err error)
	InsertMemberAccountByStruct(conn database.DBConn, ctx context.Context, ma *model.MemberAccount) error
}
