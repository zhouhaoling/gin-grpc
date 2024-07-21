package project

import (
	"context"
	"time"

	"github.com/jinzhu/copier"

	"test.com/common/errs"

	"test.com/project-api/config"
	"test.com/project-api/internal/model"
	ag "test.com/project-grpc/account_grpc"

	"github.com/gin-gonic/gin"
	"test.com/common"
)

type HandlerAccount struct {
}

func NewHandlerAccount() *HandlerAccount {
	return &HandlerAccount{}
}

func (a *HandlerAccount) account(c *gin.Context) {
	//接收请求参数
	//调用grpc project查询账号列表
	//返回数据
	//任务列表code
	var req *model.AccountReq
	_ = c.ShouldBind(&req)
	mid := c.GetInt64(config.CtxMemberIDKey)
	res := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &ag.AccountRequest{
		MemberId:         mid,
		OrganizationCode: c.GetString(config.CtxOrganizationIDKey),
		Page:             int64(req.Page),
		PageSize:         int64(req.PageSize),
		SearchType:       int32(req.SearchType),
		DepartmentCode:   req.DepartmentCode,
	}
	response, err := AccountServiceClient.Account(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var list []*model.MemberAccount
	copier.Copy(&list, response.AccountList)
	if list == nil {
		list = []*model.MemberAccount{}
	}
	var authList []*model.ProjectAuth
	copier.Copy(&authList, response.AuthList)
	if list == nil {
		authList = []*model.ProjectAuth{}
	}
	res.ResponseSuccess(c, gin.H{
		"total":    response.Total,
		"page":     req.Page,
		"list":     list,
		"authList": authList,
	})
}
