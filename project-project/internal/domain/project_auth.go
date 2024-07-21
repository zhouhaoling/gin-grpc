package domain

import (
	"context"
	"strconv"
	"time"

	"test.com/project-project/config"

	"test.com/project-project/internal/repository/database"

	"go.uber.org/zap"
	"test.com/common/errs"
	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository/dao"
)

type ProjectAuthDomain struct {
	userRpcDomain         *UserRpcDomain
	parepo                repo.ProjectAuthRepo
	projectNodeDomain     *ProjectNodeDomain
	projectAuthNodeDomain *ProjectAuthNodeDomain
	accountDomain         *AccountDomain
}

func NewProjectAuthDomain() *ProjectAuthDomain {
	return &ProjectAuthDomain{
		userRpcDomain:         NewUserRpcDomain(),
		parepo:                dao.NewProjectAuthDao(),
		projectNodeDomain:     NewProjectNodeDomain(),
		projectAuthNodeDomain: NewProjectAuthNodeDomain(),
		accountDomain:         NewAccountDomain(),
	}
}

func (d *ProjectAuthDomain) AuthList(orgCode int64) ([]*model.ProjectAuthDisplay, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	list, err := d.parepo.FindAuthList(c, orgCode)
	if err != nil {
		zap.L().Error("project AuthList projectAuthRepo.FindAuthList error", zap.Error(err))
		return nil, model.MySQLError
	}
	var pdList []*model.ProjectAuthDisplay
	for _, v := range list {
		display := v.ToDisplay()
		pdList = append(pdList, display)
	}
	return pdList, nil
}

func (d *ProjectAuthDomain) AuthListPage(orgCode, page, pageSize int64) ([]*model.ProjectAuthDisplay, int64, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	list, total, err := d.parepo.FindAuthListPage(c, orgCode, page, pageSize)
	if err != nil {
		zap.L().Error("project AuthListPage projectAuthRepo.FindAuthListPage() error", zap.Error(err))
		return nil, 0, model.MySQLError
	}
	var pdList []*model.ProjectAuthDisplay
	for _, v := range list {
		display := v.ToDisplay()
		pdList = append(pdList, display)
	}
	return pdList, total, nil
}

func (d *ProjectAuthDomain) AllNodeAndAuth(authId int64) ([]*model.ProjectNodeAuthTree, []string, *errs.BError) {
	treeList, err := d.projectNodeDomain.AllNodeList()
	if err != nil {
		zap.L().Error("project AllNodeAndAuth projectNodeDomain.AllNodeList() error", zap.Error(err))
		return nil, nil, err
	}
	authNodeList, dbErr := d.projectAuthNodeDomain.AuthNodeList(authId)
	if dbErr != nil {
		zap.L().Error("project AllNodeAndAuth projectNodeDomain.AllNodeList() error")
		return nil, nil, err
	}
	list := model.ToAuthNodeTreeList(treeList, authNodeList)
	return list, authNodeList, nil
}

func (d *ProjectAuthDomain) Save(conn database.DBConn, authId int64, nodes []string) *errs.BError {
	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeOut)
	defer cancel()
	err := d.projectAuthNodeDomain.Save(ctx, conn, authId, nodes)
	if err != nil {
		return err
	}
	return nil
}

func (d *ProjectAuthDomain) AuthNodes(mid int64) ([]string, *errs.BError) {
	account, err := d.accountDomain.FindAccount(mid)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, model.ParamsError
	}
	authorize := account.Authorize
	authId, _ := strconv.ParseInt(authorize, 10, 64)
	authNodeList, dbErr := d.projectAuthNodeDomain.AuthNodeList(authId)
	if dbErr != nil {
		return nil, model.MySQLError
	}
	return authNodeList, nil
}

func (d *ProjectAuthDomain) CreateProjectAuth(ctx context.Context, conn database.DBConn, oid int64) (*model.ProjectAuth, error) {
	pa := &model.ProjectAuth{
		Title:            "成员",
		Status:           config.Status,
		Sort:             config.Member,
		Desc:             "成员",
		CreateBy:         0,
		CreateAt:         time.Now().UnixMilli(),
		OrganizationCode: oid,
		IsDefault:        config.IsDefaultMember,
		Type:             "member",
	}
	err := d.parepo.InsertProjectAuth(ctx, conn, pa)
	if err != nil {
		zap.L().Error("project CreateProjectAuth projectNodeDomain.InsertProjectAuth() error")
		return nil, err
	}
	return pa, nil
}
