package project

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"test.com/common"
	"test.com/common/errs"
	"test.com/project-api/config"
	"test.com/project-api/internal/model"
	auth "test.com/project-grpc/auth_grpc"
)

type HandlerAuth struct {
}

func (a *HandlerAuth) authList(c *gin.Context) {
	res := common.NewResponseData()
	orgCode := c.GetString(config.CtxOrganizationIDKey)
	var page = &model.PageStruct{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeOut)
	defer cancel()

	msg := &auth.AuthRequest{
		OrganizationCode: orgCode,
		Page:             page.Page,
		PageSize:         page.PageSize,
	}
	response, err := authServiceClient.AuthList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var authList []*model.ProjectAuth
	copier.Copy(&authList, response.List)
	if authList == nil {
		authList = []*model.ProjectAuth{}
	}
	res.ResponseSuccess(c, gin.H{
		"total": response.Total,
		"list":  authList,
		"page":  page.Page,
	})

}

func (a *HandlerAuth) authApply(c *gin.Context) {
	result := common.NewResponseData()
	var req *model.ProjectAuthReq
	_ = c.ShouldBind(&req)
	var nodes []string
	if req.Nodes != "" {
		json.Unmarshal([]byte(req.Nodes), &nodes)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &auth.AuthRequest{
		Action: req.Action,
		AuthId: req.Id,
		Nodes:  nodes,
	}
	applyResponse, err := authServiceClient.AuthApply(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		result.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var list []*model.ProjectNodeAuthTree
	copier.Copy(&list, applyResponse.List)
	var checkedList []string
	copier.Copy(&checkedList, applyResponse.CheckedList)
	result.ResponseSuccess(c, gin.H{
		"list":        list,
		"checkedList": checkedList,
	})
}

func (a *HandlerAuth) GetAuthNodes(c *gin.Context) ([]string, error) {
	memberId := c.GetInt64(config.CtxMemberIDKey)
	msg := &auth.AuthRequest{
		MemberId: memberId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := authServiceClient.AuthNodesByMemberId(ctx, msg)
	if err != nil {
		return nil, err
	}
	return response.List, err
}

func NewHandlerAuth() *HandlerAuth {
	return &HandlerAuth{}
}
