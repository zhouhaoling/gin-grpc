package project

import (
	"context"
	"time"

	"test.com/common/tms"

	"test.com/project-api/config"

	"github.com/jinzhu/copier"
	"test.com/common/errs"
	pg "test.com/project-grpc/project_grpc"

	"github.com/gin-gonic/gin"
	"test.com/common"
	"test.com/project-api/internal/model"
)

type HandlerTask struct {
}

// 任务看板
func (t *HandlerTask) taskStages(c *gin.Context) {
	//1.获取参数
	//2.调用grpc服务
	//3.处理响应
	res := common.NewResponseData()
	projectCode := c.PostForm("projectCode")
	page := &model.PageStruct{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	tr := &pg.TaskRequest{
		MemberId:    c.GetInt64(config.CtxMemberIDKey),
		ProjectCode: projectCode,
		Page:        page.Page,
		PageSize:    page.PageSize,
	}
	stages, err := taskServiceClient.TaskStages(ctx, tr)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var list []*model.TaskStagesResp
	copier.Copy(&list, stages.List)
	if list == nil {
		list = []*model.TaskStagesResp{}
	}
	for _, v := range list {
		v.TasksLoading = true  //任务加载状态
		v.FixedCreator = false //添加任务按钮定位
		v.ShowTaskCard = false //是否显示创建卡片
		v.Tasks = []int{}
		v.DoneTasks = []int{}
		v.UnDoneTasks = []int{}
	}
	res.ResponseSuccess(c, gin.H{
		"list":  list,
		"total": stages.Total,
		"page":  page.Page,
	})
}

// taskMemberList 任务用户
func (t *HandlerTask) taskMemberList(c *gin.Context) {
	res := common.NewResponseData()
	projectCode := c.PostForm("projectCode")
	page := &model.PageStruct{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.TaskRequest{
		MemberId:    c.GetInt64(config.CtxMemberIDKey),
		ProjectCode: projectCode,
		Page:        page.Page,
		PageSize:    page.PageSize,
	}
	resp, err := taskServiceClient.MemberProjectList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}

	var list []*model.MemberProjectResp
	copier.Copy(&list, resp.List)
	if list == nil {
		list = []*model.MemberProjectResp{}
	}

	res.ResponseSuccess(c, gin.H{
		"list":  list,
		"total": resp.Total,
		"page":  page.Page,
	})
}

// taskList 任务列表
func (t *HandlerTask) taskList(c *gin.Context) {
	res := common.NewResponseData()
	stageCode := c.PostForm("stageCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	mid := c.GetInt64(config.CtxMemberIDKey)
	resp, err := taskServiceClient.TaskList(ctx, &pg.TaskRequest{StageCode: stageCode, MemberId: mid})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var taskDisplayList []*model.TaskDisplay
	copier.Copy(&taskDisplayList, resp.List)
	if taskDisplayList == nil {
		taskDisplayList = []*model.TaskDisplay{}
	}
	//返回给前端的数据一定不要给null
	for _, v := range taskDisplayList {
		if v.Tags == nil {
			v.Tags = []int{}
		}
		if v.ChildCount == nil {
			v.ChildCount = []int{}
		}
	}
	res.ResponseSuccess(c, taskDisplayList)
}

// taskSave 创建任务
func (t *HandlerTask) taskSave(c *gin.Context) {
	res := common.NewResponseData()
	var req *model.TaskSaveReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.TaskRequest{
		ProjectCode: req.ProjectCode,
		Name:        req.Name,
		StageCode:   req.StageCode,
		AssignTo:    req.AssignTo,
		MemberId:    c.GetInt64(config.CtxMemberIDKey),
	}
	taskMessage, err := taskServiceClient.SaveTask(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	td := &model.TaskDisplay{}
	copier.Copy(td, taskMessage)
	if td != nil {
		if td.Tags == nil {
			td.Tags = []int{}
		}
		if td.ChildCount == nil {
			td.ChildCount = []int{}
		}
	}
	res.ResponseSuccess(c, td)
}

// taskSort 移动任务
func (t *HandlerTask) taskSort(c *gin.Context) {
	res := common.NewResponseData()
	var req *model.TaskSortReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.TaskRequest{
		PreTaskCode:  req.PreTaskCode,
		NextTaskCode: req.NextTaskCode,
		ToStageCode:  req.ToStageCode,
	}
	_, err := taskServiceClient.TaskSort(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	res.ResponseSuccess(c, []int{})
}

// myTaskList 我的任务列表
func (t *HandlerTask) myTaskList(c *gin.Context) {
	res := common.NewResponseData()
	var req *model.MyTaskReq
	c.ShouldBind(&req)
	memberId := c.GetInt64(config.CtxMemberIDKey)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.TaskRequest{
		MemberId: memberId,
		TaskType: int32(req.TaskType),
		Type:     int32(req.Type),
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	myTaskListResponse, err := taskServiceClient.MyTaskList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var myTaskList []*model.MyTaskDisplay
	copier.Copy(&myTaskList, myTaskListResponse.List)
	if myTaskList == nil {
		myTaskList = []*model.MyTaskDisplay{}
	}
	for _, v := range myTaskList {
		v.ProjectInfo = model.ProjectInfo{
			Name: v.ProjectName,
			Code: v.ProjectCode,
		}
	}

	res.ResponseSuccess(c, gin.H{
		"list":  myTaskList,
		"total": myTaskListResponse.Total,
	})
}

// taskRead 任务详情
func (t *HandlerTask) taskRead(c *gin.Context) {
	res := common.NewResponseData()
	taskCode := c.PostForm("taskCode")
	mid := c.GetInt64(config.CtxMemberIDKey)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.TaskRequest{
		TaskCode: taskCode,
		MemberId: mid,
	}
	taskMessage, err := taskServiceClient.ReadTask(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	td := &model.TaskDisplay{}
	copier.Copy(td, taskMessage)
	if td != nil {
		if td.Tags == nil {
			td.Tags = []int{}
		}
		if td.ChildCount == nil {
			td.ChildCount = []int{}
		}
	}
	res.ResponseSuccess(c, td)
}

// listTaskMember 任务成员
func (t *HandlerTask) listTaskMember(c *gin.Context) {
	res := common.NewResponseData()
	taskCode := c.PostForm("taskCode")
	mid := c.GetInt64(config.CtxMemberIDKey)
	page := &model.PageStruct{}
	page.Bind(c)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.TaskRequest{
		TaskCode: taskCode,
		MemberId: mid,
		Page:     page.Page,
		PageSize: page.PageSize,
	}
	taskMemberResponse, err := taskServiceClient.ListTaskMember(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var tms []*model.TaskMember
	copier.Copy(&tms, taskMemberResponse.List)
	if tms == nil {
		tms = []*model.TaskMember{}
	}
	res.ResponseSuccess(c, gin.H{
		"list":  tms,
		"total": taskMemberResponse.Total,
		"page":  page.Page,
	})
}

// taskLog 任务日志
func (t *HandlerTask) taskLog(c *gin.Context) {
	res := common.NewResponseData()
	var req *model.TaskLogReq
	c.ShouldBind(&req)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	mid := c.GetInt64(config.CtxMemberIDKey)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	msg := &pg.TaskRequest{
		TaskCode: req.TaskCode,
		MemberId: mid,
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		All:      int32(req.All),
		Comment:  int32(req.Comment),
	}
	taskLogResponse, err := taskServiceClient.TaskLog(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var tms []*model.ProjectLogDisplay
	copier.Copy(&tms, taskLogResponse.List)
	if tms == nil {
		tms = []*model.ProjectLogDisplay{}
	}
	res.ResponseSuccess(c, gin.H{
		"list":  tms,
		"total": taskLogResponse.Total,
		"page":  req.Page,
	})
}

// taskWorkTimeList 任务工时列表
func (t *HandlerTask) taskWorkTimeList(c *gin.Context) {
	res := common.NewResponseData()
	taskCode := c.PostForm("taskCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.TaskRequest{
		TaskCode: taskCode,
		MemberId: c.GetInt64(config.CtxMemberIDKey),
	}
	taskWorkTimeResponse, err := taskServiceClient.TaskWorkTimeList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	var tms []*model.TaskWorkTime
	copier.Copy(&tms, taskWorkTimeResponse.List)
	if tms == nil {
		tms = []*model.TaskWorkTime{}
	}
	res.ResponseSuccess(c, tms)
}

// saveTaskWorkTime 保存任务工时
func (t *HandlerTask) saveTaskWorkTime(c *gin.Context) {
	res := common.NewResponseData()
	var req *model.SaveTaskWorkTimeReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.TaskRequest{
		TaskCode:  req.TaskCode,
		MemberId:  c.GetInt64(config.CtxMemberIDKey),
		Content:   req.Content,
		Num:       int32(req.Num),
		BeginTime: tms.ParseTime(req.BeginTime),
	}
	_, err := taskServiceClient.SaveTaskWorkTime(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	res.ResponseSuccess(c, []int{})
}

// taskSources 关联文件
func (t *HandlerTask) taskSources(c *gin.Context) {
	res := common.NewResponseData()
	taskCode := c.PostForm("taskCode")
	//fmt.Print("taskCode:", taskCode)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	sources, err := taskServiceClient.TaskSources(ctx, &pg.TaskRequest{TaskCode: taskCode})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}

	var slList []*model.SourceLink
	copier.Copy(&slList, sources.List)
	if slList == nil {
		slList = []*model.SourceLink{}
	}
	res.ResponseSuccess(c, slList)
}

// taskCreateComment 创建评论
func (t *HandlerTask) taskCreateComment(c *gin.Context) {
	res := common.NewResponseData()
	req := model.CommentReq{}
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &pg.TaskRequest{
		TaskCode:       req.TaskCode,
		CommentContent: req.Comment,
		Mentions:       req.Mentions,
		MemberId:       c.GetInt64(config.CtxMemberIDKey),
	}
	_, err := taskServiceClient.CreateTaskComment(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	res.ResponseSuccess(c, true)
}

func (t *HandlerTask) taskStagesSave(c *gin.Context) {
	res := common.NewResponseData()
	req := model.TaskStagesSaveReq{}
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	//fmt.Println(req)
	msg := &pg.TaskRequest{
		ProjectCode: req.ProjectCode,
		Name:        req.Name,
	}
	ts, err := taskServiceClient.CreateTaskList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	list := &model.TaskStagesResp{
		Name:         ts.Name,
		ProjectCode:  ts.ProjectCode,
		Sort:         int(ts.Sort),
		Description:  ts.Description,
		CreateTime:   ts.CreateTime,
		Code:         ts.Code,
		Deleted:      int(ts.Deleted),
		TasksLoading: true,
		FixedCreator: false,
		ShowTaskCard: false,
		Tasks:        []int{},
		DoneTasks:    []int{},
		UnDoneTasks:  []int{},
	}
	res.ResponseSuccess(c, gin.H{
		"list": list,
	})
}

func (t *HandlerTask) taskStagesDelete(c *gin.Context) {
	//任务列表code
	code := c.PostForm("code")
	res := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := taskServiceClient.DeleteTaskStages(ctx, &pg.TaskRequest{ToStageCode: code})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	res.ResponseSuccess(c, nil)
}

func (t *HandlerTask) taskStagesEdit(c *gin.Context) {
	//任务列表code
	code := c.PostForm("code")
	name := c.PostForm("name")
	res := common.NewResponseData()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := taskServiceClient.EditTaskStages(ctx, &pg.TaskRequest{ToStageCode: code, Name: name})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	res.ResponseSuccess(c, nil)
}

func NewHandlerTask() *HandlerTask {
	return &HandlerTask{}
}
