package project

import (
	"context"
	"time"

	"test.com/project-api/config"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"test.com/common"
	"test.com/common/errs"
	"test.com/project-api/internal/model"
	department "test.com/project-grpc/department_grpc"
)

type HandlerDepartment struct {
}

func NewHandlerDepartment() *HandlerDepartment {
	return &HandlerDepartment{}
}

func (h *HandlerDepartment) department(c *gin.Context) {
	res := common.NewResponseData()
	var req *model.DepartmentReq
	_ = c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &department.DepartmentReqMessage{
		Page:                 req.Page,
		PageSize:             req.PageSize,
		ParentDepartmentCode: req.Pcode,
		OrganizationCode:     c.GetString(config.CtxOrganizationIDKey),
	}
	listDepartmentMessage, err := departmentServiceClient.List(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var list []*model.Department
	copier.Copy(&list, listDepartmentMessage.List)
	if list == nil {
		list = []*model.Department{}
	}
	res.ResponseSuccess(c, gin.H{
		"total": listDepartmentMessage.Total,
		"page":  req.Page,
		"list":  list,
	})
}

func (h *HandlerDepartment) saveDepartment(c *gin.Context) {
	res := common.NewResponseData()
	var req *model.DepartmentReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &department.DepartmentReqMessage{
		Name:                 req.Name,
		DepartmentCode:       req.DepartmentCode,
		ParentDepartmentCode: req.ParentDepartmentCode,
		OrganizationCode:     c.GetString(config.CtxOrganizationIDKey),
	}
	departmentMessage, err := departmentServiceClient.Save(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var response = &model.Department{}
	copier.Copy(response, departmentMessage)
	res.ResponseSuccess(c, response)
}

func (h *HandlerDepartment) readDepartment(c *gin.Context) {
	res := common.NewResponseData()
	departmentCode := c.PostForm("departmentCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &department.DepartmentReqMessage{
		DepartmentCode:   departmentCode,
		OrganizationCode: c.GetString(config.CtxOrganizationIDKey),
	}
	departmentMessage, err := departmentServiceClient.ReadDepartment(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var response = &model.Department{}
	copier.Copy(response, departmentMessage)
	res.ResponseSuccess(c, response)
}
