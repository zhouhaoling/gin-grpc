package repo

import (
	"context"

	"test.com/project-user/internal/repository/database"

	"test.com/project-user/internal/model"
)

type Organization interface {
	InsertOrganization(conn database.DBConn, ctx context.Context, org *model.Organization) error
	SelectOrganizationListByMId(ctx context.Context, mid int64) ([]*model.Organization, error)
}
