package repository

import (
	"context"
	"time"

	"test.com/project-user/internal/repository/dao"

	"test.com/project-user/config"
	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repo"
	"test.com/project-user/internal/repository/database"
)

// OrganizationRepository 组织的数据库操作
type OrganizationRepository struct {
	dao repo.OrganizationRepo
}

func NewOrganizationRepository() *OrganizationRepository {
	return &OrganizationRepository{
		dao: dao.NewOrganizationDao(),
	}
}

func (repo *OrganizationRepository) CreateOrganization(conn database.DBConn, ctx context.Context, member *model.Member) (*model.Organization, error) {
	org := &model.Organization{
		Name:       member.Name + config.PersonalOrganization,
		MemberId:   member.MId,
		CreateTime: time.Now().UnixMilli(),
		Personal:   config.Personal,
		//存储的是图片的路径
		Avatar: config.DefaultAvatar,
	}
	return repo.dao.InsertOrganization(conn, ctx, org)
}

func (repo *OrganizationRepository) FindOrganizationListByMId(ctx context.Context, mid int64) ([]*model.Organization, error) {
	return repo.dao.SelectOrganizationListByMId(ctx, mid)
}
