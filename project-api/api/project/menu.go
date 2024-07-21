package project

import (
	"context"

	mg "test.com/project-grpc/menu_grpc"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"test.com/common"
	"test.com/common/errs"
	"test.com/project-api/config"
	"test.com/project-api/internal/model"
)

type HandlerMenu struct {
}

func (d *HandlerMenu) menuList(c *gin.Context) {
	result := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeOut)
	defer cancel()
	res, err := menuServiceClient.MenuList(ctx, &mg.MenuReqMessage{})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		result.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var list []*model.Menu
	copier.Copy(&list, res.List)
	if list == nil {
		list = []*model.Menu{}
	}
	result.ResponseSuccess(c, list)
}

func NewHandlerMenu() *HandlerMenu {
	return &HandlerMenu{}
}
