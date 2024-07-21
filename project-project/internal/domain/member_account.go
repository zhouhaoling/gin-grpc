package domain

import (
	"context"
	"fmt"
	"time"

	"test.com/project-project/internal/repository/database"

	"go.uber.org/zap"

	"test.com/common/encrypts"
	"test.com/project-project/internal/model"

	"test.com/common/errs"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository/dao"
)

type AccountDomain struct {
	arepo            repo.AccountRepo
	departmentDomain *DepartmentDomain
	userRpcDomain    *UserRpcDomain
}

func (d *AccountDomain) AccountList(organizationCode string, memberId int64, page int64, pageSize int64, departmentCode string, searchType int32) ([]*model.MemberAccountDisplay, int64, *errs.BError) {
	condition := ""
	organizationCodeId := encrypts.DecryptInt64(organizationCode)
	departmentCodeId := encrypts.DecryptInt64(departmentCode)
	switch searchType {
	case 1:
		condition = "status = 1"
	case 2:
		condition = "department_code = NULL"
	case 3:
		condition = "status = 0"
	case 4:
		condition = fmt.Sprintf("status = 1 and department_code = %d", departmentCodeId)
	default:
		condition = "status = 1"
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	list, total, err := d.arepo.FindList(c, condition, organizationCodeId, departmentCodeId, page, pageSize)
	if err != nil {
		zap.L().Error("AccountDomain AccountList FindList() failed", zap.Error(err))
		return nil, 0, model.MySQLError
	}
	//fmt.Println(*&list)
	var dList []*model.MemberAccountDisplay
	for _, v := range list {
		display := v.ToDisplay()
		//TODO查询用户获取头像信息，或者直接用ms_account_member里面存储的头像信息
		//memberInfo, _ := d.userRpcDomain.MemberInfo(v.MemberCode)
		//display.Avatar = memberInfo.Avatar
		if v.DepartmentCode > 0 {
			department, err := d.departmentDomain.FindDepartmentById(v.DepartmentCode)
			if err != nil {
				return nil, 0, err
			}
			display.Departments = department.Name
		}
		//fmt.Printf("%+v\n", display)
		dList = append(dList, display)
	}
	return dList, total, nil
}

func (d *AccountDomain) CreateMemberAccount(ctx context.Context, conn database.DBConn, ma *model.MemberAccount) error {
	err := d.arepo.InsertMemberAccountByStruct(conn, ctx, ma)
	if err != nil {
		zap.L().Error("domain member_account CreateMemberAccount arepo.InsertMemberAccountByStruct()", zap.Error(err))
		return err
	}
	return nil
}

func (d *AccountDomain) FindAccount(mid int64) (*model.MemberAccount, *errs.BError) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	account, err := d.arepo.SelectAccountByMemberId(ctx, mid)
	if err != nil {
		zap.L().Error("AccountDomain AccountList SelectAccountByMemberId() failed", zap.Error(err))
		return nil, model.MySQLError
	}
	return account, nil
}

func NewAccountDomain() *AccountDomain {
	return &AccountDomain{
		userRpcDomain:    NewUserRpcDomain(),
		arepo:            dao.NewAccountDao(),
		departmentDomain: NewDepartmentDomain(),
	}
}
