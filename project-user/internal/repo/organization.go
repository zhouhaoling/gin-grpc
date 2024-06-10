package repo

import (
	"context"

	"test.com/project-user/internal/repository/database"

	"test.com/project-user/internal/model"
)

type Organization interface {
	InsertOrganization(conn database.DBConn, ctx context.Context, org *model.Organization) error
}
