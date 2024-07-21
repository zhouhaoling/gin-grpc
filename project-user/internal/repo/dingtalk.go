package repo

import (
	"context"
	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/database"
)

type DingTalkRepo interface {
	InsertDingTalk(conn database.DBConn, ctx context.Context, dingTalk *model.DingTalk) error
}
