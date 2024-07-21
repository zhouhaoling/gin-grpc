package project

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"test.com/common"
	"test.com/common/errs"
	"test.com/project-api/config"
	"test.com/project-api/internal/model"
	pg "test.com/project-grpc/project_grpc"
)

type HandlerProject struct {
}

func NewHandlerProject() *HandlerProject {
	return &HandlerProject{}
}

// index 首页
func (p *HandlerProject) index(c *gin.Context) {
	//1.调用grpc服务查询数据
	//2.返回响应
	res := &common.ResponseData{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.IndexMessage{}
	indexResponse, err := projectServiceClient.Index(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	menus := indexResponse.Menus
	var menu []*model.Menu
	copier.Copy(&menu, menus)
	fmt.Println("menu:", menu)
	res.ResponseSuccess(c, menu)
}

// myProjectList 我的项目列表
func (p *HandlerProject) myProjectList(c *gin.Context) {
	res := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	mid := c.GetInt64(config.CtxMemberIDKey)
	name := c.GetString(config.CtxMemberNameKey)
	page := &model.PageStruct{}
	page.Bind(c)
	selectBy := c.PostForm("selectBy")
	pr := &pg.ProjectRequest{
		MemberId:   mid,
		Page:       page.Page,
		PageSize:   page.PageSize,
		MemberName: name,
		SelectBy:   selectBy,
	}
	resp, err := projectServiceClient.FindProjectByMemberId(ctx, pr)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}

	var pam []*model.ProjectAndMember
	err = copier.Copy(&pam, resp.Pm)
	if err != nil {
		zap.L().Error("copier.Copy(&pam, resp.Pm) failed", zap.Error(err))
		res.ResponseError(c, common.CodeServerBusy)
		return
	}
	if pam == nil {
		pam = []*model.ProjectAndMember{}
	}
	res.ResponseSuccess(c, gin.H{
		"list":  pam,
		"total": resp.Total,
	})
}

// projectTemplate 项目模板
func (p *HandlerProject) projectTemplate(c *gin.Context) {
	res := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	mid := c.GetInt64(config.CtxMemberIDKey)
	name := c.GetString(config.CtxMemberNameKey)
	orgCode := c.GetString(config.CtxOrganizationIDKey)
	var page = &model.PageStruct{}
	page.Bind(c)
	viewTypeStr := c.PostForm("viewType")
	viewType, _ := strconv.ParseInt(viewTypeStr, 10, 64)
	pm := &pg.ProjectRequest{
		MemberId:         mid,
		MemberName:       name,
		OrganizationCode: orgCode,
		Page:             page.Page,
		PageSize:         page.PageSize,
		ViewType:         int32(viewType),
	}
	templateResp, err := projectServiceClient.FindProjectTemplate(ctx, pm)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
	}
	var pts []*model.ProjectTemplate
	copier.Copy(&pts, templateResp.Ptm)
	if pts == nil {
		pts = []*model.ProjectTemplate{}
	}
	for _, v := range pts {
		if v.TaskStages == nil {
			v.TaskStages = []*model.TaskStagesOnlyName{}
		}
	}
	res.ResponseSuccess(c, gin.H{
		"list":  pts,
		"total": templateResp.Total,
	})
}

// createProject 创建项目
func (p *HandlerProject) createProject(c *gin.Context) {
	//获取请求参数
	//绑定参数,调用grpc的创建项目服务
	//返回响应
	res := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	mid := c.GetInt64(config.CtxMemberIDKey)
	orgCode := c.GetString(config.CtxOrganizationIDKey)
	var param model.ParamProject
	err := param.Bind(c)
	if err != nil {
		zap.L().Error("param.Bind(c) failed", zap.Error(err))
		res.ResponseError(c, common.CodeInvalidParams)
		return
	}
	pt := &pg.ProjectRequest{
		MemberId:         mid,
		OrganizationCode: orgCode,
		Name:             param.Name,
		TemplateCode:     param.TemplateCode,
		Description:      param.Description,
		Id:               int64(param.Id),
	}
	cpResp, err := projectServiceClient.CreateProject(ctx, pt)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var rsp *model.SaveProject
	copier.Copy(&rsp, cpResp)
	res.ResponseSuccess(c, rsp)
	return
}

// readProject 读取项目详情
func (p *HandlerProject) readProject(c *gin.Context) {
	res := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	pCodeStr := c.PostForm("projectCode")
	//fmt.Println("pCodeStr:", pCodeStr)
	memberId := c.GetInt64(config.CtxMemberIDKey)
	pr := &pg.ProjectRequest{
		MemberId:    memberId,
		ProjectCode: pCodeStr,
	}
	resp, err := projectServiceClient.FindProjectDetail(ctx, pr)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var pd = &model.ProjectDetail{}
	err = copier.Copy(pd, resp)
	if err != nil {
		zap.L().Error("copier.Copy(&pd, resp) failed", zap.Error(err))
		res.ResponseError(c, common.CodeServerBusy)
		return
	}
	res.ResponseSuccess(c, pd)
}

// recycleProject 回收项目
func (p *HandlerProject) recycleProject(c *gin.Context) {
	res := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	pCodeStr := c.PostForm("projectCode")
	_, err := projectServiceClient.RecycleProjectByPid(ctx, &pg.ProjectRequest{ProjectCode: pCodeStr, Deleted: true})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	res.ResponseSuccess(c, []int{})
}

// recoveryProject 恢复项目
func (p *HandlerProject) recoveryProject(c *gin.Context) {
	res := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	pCodeStr := c.PostForm("projectCode")
	_, err := projectServiceClient.RecycleProjectByPid(ctx, &pg.ProjectRequest{ProjectCode: pCodeStr, Deleted: false})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	res.ResponseSuccess(c, []int{})
}

// collectProject 收藏项目
func (p *HandlerProject) collectProject(c *gin.Context) {
	res := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	pCodeStr := c.PostForm("projectCode")
	mid := c.GetInt64(config.CtxMemberIDKey)
	collectType := c.PostForm("type")
	pr := &pg.ProjectRequest{
		MemberId:    mid,
		ProjectCode: pCodeStr,
		CollectType: collectType,
	}
	_, err := projectServiceClient.CollectProjectByType(ctx, pr)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	if collectType == "cancel" {
		res.ResponseSuccessWithMsg(c, "取消收藏成功", []int{})
		return
	}
	res.ResponseSuccessWithMsg(c, "加入收藏成功", []int{})
}

// editProject 编辑项目
func (p *HandlerProject) editProject(c *gin.Context) {
	//获取请求参数
	//绑定参数,调用grpc的编写项目服务
	//返回响应
	res := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	var param model.ParamProjectEdit
	err := param.Bind(c)
	if err != nil {
		zap.L().Error("param.Bind(c) failed", zap.Error(err))
		res.ResponseError(c, common.CodeInvalidParams)
		return
	}
	mid := c.GetInt64(config.CtxMemberIDKey)
	per := &pg.ProjectEditRequest{
		ProjectCode:        param.ProjectCode,
		Cover:              param.Cover,
		Name:               param.Name,
		Description:        param.Description,
		Schedule:           float64(param.Schedule),
		Private:            int32(param.Private),
		Prefix:             param.Prefix,
		OpenPrefix:         int32(param.OpenPrefix),
		OpenBeginTime:      int32(param.OpenBeginTime),
		OpenTaskPrivate:    int32(param.OpenTaskPrivate),
		TaskBoardTheme:     param.TaskBoardTheme,
		AutoUpdateSchedule: int32(param.AutoUpdateSchedule),
		MemberCode:         mid,
	}
	fmt.Println("pCode", per.ProjectCode)
	_, err = projectServiceClient.EditProject(ctx, per)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	res.ResponseSuccess(c, []int{})
}

func (p *HandlerProject) getLogBySelfProject(c *gin.Context) {
	res := common.NewResponseData()
	var page = &model.PageStruct{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.ProjectRequest{
		MemberId: c.GetInt64(config.CtxMemberIDKey),
		Page:     page.Page,
		PageSize: page.PageSize,
	}
	projectLogResponse, err := projectServiceClient.GetLogBySelfProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var list []*model.ProjectLog
	copier.Copy(&list, projectLogResponse.List)
	if list == nil {
		list = []*model.ProjectLog{}
	}
	res.ResponseSuccess(c, list)
}

func (p *HandlerProject) nodeList(c *gin.Context) {
	result := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	response, err := projectServiceClient.NodeList(ctx, &pg.ProjectRequest{})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		result.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var list []*model.ProjectNodeTree
	copier.Copy(&list, response.Nodes)

	result.ResponseSuccess(c, gin.H{
		"nodes": list,
	})
}

func (p *HandlerProject) FindProjectByMemberId(mid int64, projectCode string, taskCode string) (*model.Project, bool, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.ProjectRequest{
		MemberId:    mid,
		ProjectCode: projectCode,
		TaskCode:    taskCode,
	}
	projectResponse, err := projectServiceClient.FindProjectByMid(ctx, msg)
	if err != nil {
		return nil, false, false, err
	}
	if projectResponse.Project == nil {
		return nil, false, false, nil
	}
	pr := &model.Project{}
	copier.Copy(pr, projectResponse.Project)
	return pr, true, projectResponse.IsOwner, nil
}
