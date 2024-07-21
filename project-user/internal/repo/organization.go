package repo

import (
	"context"

	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/database"
)

type OrganizationRepo interface {
	InsertOrganization(conn database.DBConn, ctx context.Context, org *model.Organization) (*model.Organization, error)
	SelectOrganizationListByMId(ctx context.Context, mid int64) ([]*model.Organization, error)
	SelectOrganizationByOrgId(ctx context.Context, id int64) (bool, error)
}
