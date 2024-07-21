package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"test.com/project-project/internal/domain"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"test.com/common/encrypts"
	"test.com/common/errs"
	"test.com/common/tms"
	pg "test.com/project-grpc/project_grpc"
	"test.com/project-project/config"
	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository"
	"test.com/project-project/internal/repository/dao"
	"test.com/project-project/internal/repository/database"
)

type ProjectService struct {
	pg.UnimplementedProjectServiceServer
	cache                 repo.Cache
	tran                  database.Transaction
	mrepo                 *repository.MenuRepository
	prepo                 *repository.ProjectRepository
	trepo                 *repository.ProjectTemplateRepository
	tsrepo                *repository.TaskStagesRepository
	udomain               *domain.UserRpcDomain
	nodeDomain            *domain.ProjectNodeDomain
	taskDomain            *domain.TaskDomain
	projectAuthDomain     *domain.ProjectAuthDomain
	projectAuthNodeDomain *domain.ProjectAuthNodeDomain
	accountDomain         *domain.AccountDomain
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		cache:                 dao.NewRedisCache(),
		tran:                  repository.NewTransaction(),
		mrepo:                 repository.NewMenuRepository(),
		prepo:                 repository.NewProjectRepository(),
		trepo:                 repository.NewProjectTemplateRepository(),
		tsrepo:                repository.NewTaskStagesRepository(),
		udomain:               domain.NewUserRpcDomain(),
		nodeDomain:            domain.NewProjectNodeDomain(),
		taskDomain:            domain.NewTaskDomain(),
		projectAuthDomain:     domain.NewProjectAuthDomain(),
		projectAuthNodeDomain: domain.NewProjectAuthNodeDomain(),
		accountDomain:         domain.NewAccountDomain(),
	}
}

// Index 首页
func (svc *ProjectService) Index(ctx context.Context, msg *pg.IndexMessage) (*pg.IndexResponse, error) {
	pms, err := svc.mrepo.FindMenu(context.Background())
	//pms, err := svc.menuRepo.SelectMenus(context.Background())
	if err != nil {
		zap.L().Error("find menu failed, error:", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	childs := model.CovertChild(pms)
	var mms []*pg.MenuMessage
	copier.Copy(&mms, childs)

	return &pg.IndexResponse{Menus: mms}, nil
}

// FindProjectByMemberId 查询项目列表
func (svc *ProjectService) FindProjectByMemberId(ctx context.Context, msg *pg.ProjectRequest) (*pg.ProjectResponse, error) {
	mid := msg.MemberId
	page := msg.Page
	size := msg.PageSize
	fmt.Println("mid:", msg.MemberId, "name:", msg.MemberName)
	//对查询类型进行判断
	var pm []*model.ProjectAndMember
	var total int64
	var err error
	if msg.SelectBy == "" || msg.SelectBy == "my" {
		pm, total, err = svc.prepo.FindProjectByMemId(context.Background(), mid, page, size, 0)
	}
	if msg.SelectBy == "archive" {
		pm, total, err = svc.prepo.FindProjectByMemId(context.Background(), mid, page, size, 1)
	}
	if msg.SelectBy == "deleted" {
		pm, total, err = svc.prepo.FindProjectByMemId(context.Background(), mid, page, size, 2)

	}
	if msg.SelectBy == "collect" {
		pm, total, err = svc.prepo.FindProjectByMemId(context.Background(), mid, page, size, 3)
		for _, v := range pm {
			v.Collected = config.Collected
		}
	} else {
		collectPms, _, err := svc.prepo.FindProjectByMemId(context.Background(), mid, page, size, 3)
		if err != nil {
			zap.L().Error("project FindProjectByMemberId is failed, error:", zap.Error(err))
			return nil, errs.GrpcError(model.MySQLError)
		}
		var cMap = make(map[int64]*model.ProjectAndMember)
		for _, v := range collectPms {
			cMap[v.Id] = v
		}
		for _, v := range pm {
			if cMap[v.ProjectCode] != nil {
				v.Collected = config.Collected
			}
		}
	}

	if err != nil {
		zap.L().Error("project FindProjectByMemberId is failed, error:", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if pm == nil {
		return &pg.ProjectResponse{Pm: []*pg.ProjectMessage{}, Total: total}, nil
	}
	var pmm []*pg.ProjectMessage
	err = copier.Copy(&pmm, pm)
	if err != nil {
		zap.L().Error("copy pm => pmm failed", zap.Error(err))
		return nil, errs.GrpcError(model.ErrorServerBusy)
	}
	for _, v := range pmm {
		v.Code, err = encrypts.EncryptInt64(v.ProjectCode, config.AESKey)
		pam := model.ToMap(pm)[v.Id]
		v.AccessControlType = pam.GetAccessControlType()
		v.OrganizationCode, _ = encrypts.EncryptInt64(pam.OrganizationCode, config.AESKey)
		v.JoinTime = tms.FormatByMill(pam.JoinTime)
		v.OwnerName = msg.MemberName
		v.Order = int32(pam.Sort)
		v.CreateTime = tms.FormatByMill(pam.CreateTime)
		if msg.SelectBy == "deleted" {
			v.DeletedTime = tms.FormatByMill(pam.DeletedTime)
		}
	}
	return &pg.ProjectResponse{Pm: pmm, Total: total}, nil
}

func (svc *ProjectService) FindProjectTemplate(ctx context.Context, msg *pg.ProjectRequest) (*pg.ProjectTemplateResponse, error) {
	//根据viewType查询项目模板表得到list
	//模型转换,拿到模板id列表 去 任务模板表查询
	//组装数据
	//到模板用户的项目模板去拿信息
	//TODO 获取模板用户的项目模板去拿信息
	orgIdStr, err := encrypts.Decrypt(msg.OrganizationCode, config.AESKey)
	if err != nil {
		zap.L().Error("FindProjectTemplate Decrypt msg.OrganizationCode failed", zap.Error(err))
		return nil, errs.GrpcError(model.ErrorServerBusy)
	}
	orgId, _ := strconv.ParseInt(orgIdStr, 10, 64)
	//查询项目列表得到list
	pts, total, err := svc.trepo.FindProjectTemplate(ctx, msg, orgId)
	if err != nil {
		zap.L().Error("dao.SelectProjectTemplate() failed", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	//得到ids list ，到任务模板列表中查询
	ids := model.ToProjectTemplateIds(pts)
	mtst, err := svc.trepo.FindInProTemIds(ctx, ids)
	if err != nil {
		zap.L().Error("dao.SelectInProTemIds() failed", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	//组装数据
	var ptall []*model.ProjectTemplateAll
	for _, v := range pts {
		ptall = append(ptall, v.Convert(model.CovertProjectMap(mtst)[v.Id]))
	}
	var pmMsg []*pg.ProjectTemplateMessage
	copier.Copy(&pmMsg, ptall)

	return &pg.ProjectTemplateResponse{Ptm: pmMsg, Total: total}, nil
}

func (svc *ProjectService) CopyProjectTemplate(ctx context.Context, msg *pg.ProjectRequest) (*pg.CopyProjectTemplateResponse, error) {
	oid := encrypts.DecryptInt64(msg.OrganizationCode)
	mid := msg.MemberId
	var defaultMid int64 = 21501921651068928
	//获取模板用户的项目模板信息
	templateList, err := svc.trepo.FindProjectTemplateByMid(ctx, defaultMid)
	if err != nil {
		zap.L().Error("project CopyProjectTemplate FindProjectTemplateByMid() failed", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	newTemplateList := make([]*model.ProjectTemplate, len(templateList))
	for index, template := range templateList {
		t := &model.ProjectTemplate{}
		t.Name = template.Name
		t.OrganizationCode = oid
		t.MemberCode = mid
		t.Sort = template.Sort
		t.Description = template.Description
		t.Cover = template.Cover
		t.IsSystem = template.IsSystem
		t.CreateTime = time.Now().UnixMilli()
		newTemplateList[index] = t
	}
	//获取模板用户的项目模板id列表
	var ptcodes []int64
	for _, template := range templateList {
		ptcodes = append(ptcodes, int64(template.Id))
	}
	//根据模板用户的项目模板id查询任务阶段模板信息
	ptList, err := svc.trepo.FindTaskStagesTemplateByPtcodes(ctx, ptcodes)
	if err != nil {
		zap.L().Error("project CopyProjectTemplate FindTaskStagesTemplateByPtcodes() failed", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	newTaskStagesTemplateList := make([]*model.MsTaskStagesTemplate, len(ptList))
	for index, v := range ptList {
		t := &model.MsTaskStagesTemplate{}
		t.Name = v.Name
		t.Sort = v.Sort
		t.CreateTime = time.Now().UnixMilli()
		t.ProjectTemplateCode = v.ProjectTemplateCode
		newTaskStagesTemplateList[index] = t
	}
	//查询模板用户的节点列表
	projectAuthNodeList, err := svc.projectAuthNodeDomain.FindProjectAuthNode(ctx)
	if err != nil {
		return nil, errs.GrpcError(model.MySQLError)
	}
	//开启事物
	var projectAuth *model.ProjectAuth
	//创建用户项目模板
	err = svc.tran.Action(func(conn database.DBConn) error {
		//创建该用户的项目模板
		err = svc.trepo.CreateProjectTemplate(ctx, conn, newTemplateList)
		if err != nil {
			zap.L().Error("project CopyProjectTemplate CreateProjectTemplate() failed", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		fmt.Println("newTaskStagesTemplateList", newTaskStagesTemplateList)
		for _, v := range newTaskStagesTemplateList {
			if v.ProjectTemplateCode == 11 {
				v.ProjectTemplateCode = newTemplateList[0].Id
			}
			if v.ProjectTemplateCode == 12 {
				v.ProjectTemplateCode = newTemplateList[1].Id
			}
			if v.ProjectTemplateCode == 13 {
				v.ProjectTemplateCode = newTemplateList[2].Id
			}
			if v.ProjectTemplateCode == 19 {
				v.ProjectTemplateCode = newTemplateList[3].Id
			}
		}
		//创建该用户的任务模板
		err = svc.trepo.CreateTaskStagesTemplate(ctx, conn, newTaskStagesTemplateList)
		if err != nil {
			zap.L().Error("project CopyProjectTemplate CreateTaskStagesTemplate() failed", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		//创建该用户的ms_project_auth和ms_project_auth_node
		projectAuth, err = svc.projectAuthDomain.CreateProjectAuth(ctx, conn, oid)
		if err != nil {
			return errs.GrpcError(model.MySQLError)
		}
		err = svc.projectAuthNodeDomain.CreateProjectAuthNode(ctx, conn, projectAuth, projectAuthNodeList)
		if err != nil {
			return errs.GrpcError(model.MySQLError)
		}

		//创建用户的ms_member_account
		ma := &model.MemberAccount{
			MemberCode:       mid,
			OrganizationCode: oid,
			Authorize:        strconv.FormatInt(projectAuth.Id, 10),
			IsOwner:          config.Owner,
			Name:             msg.Name,
			Mobile:           msg.Mobile,
			Email:            msg.Email,
			CreateTime:       time.Now().UnixMilli(),
			Status:           config.Status,
			Description:      msg.Description,
			Avatar:           msg.Avatar,
		}
		err = svc.accountDomain.CreateMemberAccount(ctx, conn, ma)
		if err != nil {
			return errs.GrpcError(model.MySQLError)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pg.CopyProjectTemplateResponse{}, nil
}

func (svc *ProjectService) CreateProject(ctxs context.Context, msg *pg.ProjectRequest) (*pg.SaveProjectResponse, error) {
	orgCodeStr, _ := encrypts.Decrypt(msg.OrganizationCode, config.AESKey)
	orgCode, _ := strconv.ParseInt(orgCodeStr, 10, 64)
	tplCodeStr, _ := encrypts.Decrypt(msg.TemplateCode, config.AESKey)
	tplCode, _ := strconv.ParseInt(tplCodeStr, 10, 64)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancelFunc()
	taskTemplateList, err := svc.trepo.FindTaskTemplateByProjectTemplateId(ctx, int(tplCode))
	if err != nil {
		zap.L().Error("project CreateProject FindTaskTemplateByProjectTemplateId() failed", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	//保存项目表
	project := &model.Project{
		Name:              msg.Name,
		Cover:             config.DefaultProjectCover,
		Description:       msg.Description,
		TemplateCode:      tplCode,
		OrganizationCode:  orgCode,
		CreateTime:        time.Now().UnixMilli(),
		Deleted:           config.NoDeleted,
		Archive:           config.NoArchive,
		AccessControlType: config.Open,
		TaskBoardTheme:    config.Simple,
	}

	err = svc.tran.Action(func(conn database.DBConn) error {
		//1.保存项目表
		err := svc.prepo.CreateProject(conn, ctx, project)
		if err != nil {
			zap.L().Error("CreateProject InsertProject dao failed", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		pm := &model.ProjectMember{
			ProjectCode: project.Id,
			MemberCode:  msg.MemberId,
			JoinTime:    time.Now().UnixMilli(),
			IsOwner:     msg.MemberId,
			Authorize:   config.DefaultAuthorize,
		}
		//2.保存项目和成员的关系
		err = svc.prepo.CreateProjectMember(conn, ctx, pm)
		if err != nil {
			zap.L().Error("CreateProject CreateProjectMember dao failed", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		//3.生成任务步骤
		for index, v := range taskTemplateList {
			taskStage := &model.TaskStages{
				ProjectCode: project.Id,
				Name:        v.Name,
				Sort:        index + 1,
				Description: "",
				CreateTime:  time.Now().UnixMilli(),
				Deleted:     config.NoDeleted,
			}
			err = svc.tsrepo.CreateTaskStage(conn, ctx, taskStage)
			if err != nil {
				zap.L().Error("project CreateProject CreateTaskStage()  failed", zap.Error(err))
				return errs.GrpcError(model.MySQLError)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	code, _ := encrypts.EncryptInt64(project.Id, config.AESKey)

	spr := &pg.SaveProjectResponse{
		Id:               project.Id,
		Code:             code,
		OrganizationCode: orgCodeStr,
		Name:             project.Name,
		Cover:            project.Cover,
		Description:      project.Description,
		CreateTime:       tms.FormatByMill(project.CreateTime),
		TaskBoardTheme:   project.TaskBoardTheme,
	}
	return spr, nil
}

func (svc *ProjectService) FindProjectDetail(ctx context.Context, msg *pg.ProjectRequest) (*pg.ProjectDetailResponse, error) {
	//1. 查项目表
	//2. 项目和成员的关联表 查到项目的拥有者 去member表查名字
	//3. 查收藏表 判断收藏状态
	pCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, config.AESKey)
	fmt.Println("pCodeStr", pCodeStr)
	pCode, _ := strconv.ParseInt(pCodeStr, 10, 64)
	memberId := msg.MemberId
	c, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	//到项目表查询项目信息
	projectAndmember, err := svc.prepo.FindProjectByPidAndMemberId(c, pCode, memberId)
	if err != nil {
		if errors.Is(err, model.ErrorServerBusy) {
			return nil, err
		}
		zap.L().Error("project FindProjectDetail FindProjectByPIdAndMemId error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if projectAndmember == nil {
		return nil, errs.GrpcError(model.ErrorServerBusy)
	}
	isOwner := projectAndmember.IsOwner
	//调用user grpc服务查询用户信息
	//member, err := rpc.UserServiceClient.FindMemberByMemId(c, &pu.UserRequest{MemberId: isOwner})
	member, err := svc.udomain.MemberInfo(isOwner)
	if err != nil {
		zap.L().Error("project FindProjectDetail udomain.MemberInfo(isOwner) error", zap.Error(err))
		return nil, err
	}
	isCollect, err := svc.prepo.FindCollectByPidAndMemId(ctx, pCode, memberId)
	if err != nil {
		zap.L().Error("project FindProjectDetail FindProjectByPIdAndMemId error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	//TODO 优化收藏的时候放入redis，直接从redis中查
	if isCollect {
		projectAndmember.Collected = config.Collected
	}

	var detailMsg = &pg.ProjectDetailResponse{}
	err = copier.Copy(detailMsg, projectAndmember)
	if err != nil {
		zap.L().Error("拷贝失败", zap.Error(err))
		return nil, errs.GrpcError(model.ErrorServerBusy)
	}
	detailMsg.OwnerAvatar = member.Avatar
	detailMsg.OwnerName = member.Name
	detailMsg.Code = msg.ProjectCode
	detailMsg.AccessControlType = projectAndmember.GetAccessControlType()
	detailMsg.OrganizationCode, _ = encrypts.EncryptInt64(projectAndmember.OrganizationCode, config.AESKey)
	fmt.Println("projectAndMember.OrganizationCode:", projectAndmember.OrganizationCode)
	detailMsg.Order = int32(projectAndmember.Sort)
	detailMsg.CreateTime = tms.FormatByMill(projectAndmember.CreateTime)
	//fmt.Println("detailMsg", detailMsg)
	return detailMsg, nil
}

// RecycleProjectByPid 移入和移出回收站
func (svc *ProjectService) RecycleProjectByPid(ctx context.Context, msg *pg.ProjectRequest) (*pg.IndexMessage, error) {
	pCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, config.AESKey)
	fmt.Println("pCodeStr", pCodeStr)
	pCode, _ := strconv.ParseInt(pCodeStr, 10, 64)
	c, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err := svc.prepo.RecycleProjectByPid(c, pCode, msg.Deleted)
	if err != nil {
		zap.L().Error("project RecycleProjectByPid RecycleProjectByPid error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	return &pg.IndexMessage{}, nil
}

// CollectProjectByType 收藏或移出收藏表
func (svc *ProjectService) CollectProjectByType(ctx context.Context, msg *pg.ProjectRequest) (*pg.CollectProjectResponse, error) {
	//查项目表项目是否存在
	//根据collectType插入或者删除收藏表
	//返回响应
	pCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, config.AESKey)
	pCode, _ := strconv.ParseInt(pCodeStr, 10, 64)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := svc.prepo.CollectProjectByType(c, pCode, msg.MemberId, msg.CollectType)
	if err != nil {
		zap.L().Error("project CollectProjectByType CollectProjectByType error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	return &pg.CollectProjectResponse{}, nil
}

func (svc *ProjectService) EditProject(ctx context.Context, msg *pg.ProjectEditRequest) (*pg.ProjectEditResponse, error) {
	pCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, config.AESKey)
	pCode, _ := strconv.ParseInt(pCodeStr, 10, 64)
	//fmt.Println("pCode:", pCode)
	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	pr := &model.Project{
		Name:               msg.Name,
		Description:        msg.Description,
		Cover:              msg.Cover,
		TaskBoardTheme:     msg.TaskBoardTheme,
		Prefix:             msg.Prefix,
		Private:            int(msg.Private),
		OpenPrefix:         int(msg.OpenPrefix),
		OpenBeginTime:      int(msg.OpenBeginTime),
		OpenTaskPrivate:    int(msg.OpenTaskPrivate),
		Schedule:           msg.Schedule,
		AutoUpdateSchedule: int(msg.AutoUpdateSchedule),
	}
	err := svc.prepo.UpdateProjectByStruct(c, pCode, pr)
	if err != nil {
		zap.L().Error("project EditProject UpdateProjectByStruct error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	//TODO 将项目的头像上传到COS,并获取地址存储到数据库,未实现
	return &pg.ProjectEditResponse{}, nil
}

func (svc *ProjectService) GetLogBySelfProject(ctxs context.Context, msg *pg.ProjectRequest) (*pg.ProjectLogResponse, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc()
	//根据用户id查询当前用户的日志表信息
	logList, total, err := svc.prepo.FindProjectLogByMid(ctx, msg.MemberId, msg.Page, msg.PageSize)
	if err != nil {
		zap.L().Error("project GetLogBySelfProject FindProjectLogByMid() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	//查询项目信息
	pidList := make([]int64, len(logList))
	midList := make([]int64, len(logList))
	tidList := make([]int64, len(logList))
	for _, v := range logList {
		pidList = append(pidList, v.ProjectCode)
		midList = append(midList, v.MemberCode)
		tidList = append(tidList, v.SourceCode)
	}
	projectList, err := svc.prepo.FindProjectByPIds(ctx, pidList)
	if err != nil {
		zap.L().Error("project GetLogBySelfProject FindProjectByPIds() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	pMap := make(map[int64]*model.Project)
	for _, v := range projectList {
		pMap[v.Id] = v
	}
	//memberList, _ := rpc.UserServiceClient.FindMemberInfoByMIds(ctx, &pu.UserRequest{MIds: midList})
	//mMap := make(map[int64]*pu.MemberResponse)
	//for _, v := range memberList.List {
	//	mMap[v.Mid] = v
	//}
	_, mMap, err := svc.udomain.MemberList(midList)
	if err != nil {
		zap.L().Error("project GetLogBySelfProject udomain.MemberList(midList) error", zap.Error(err))
		return nil, err
	}
	taskList, err := svc.tsrepo.FindTaskByIds(ctx, tidList)
	if err != nil {
		zap.L().Error("project GetLogBySelfProject FindMemberInfoByMIds() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	tMap := make(map[int64]*model.Task)
	for _, v := range taskList {
		tMap[v.Id] = v
	}
	var list []*model.IndexProjectLogDisplay
	for _, v := range logList {
		display := v.ToIndexDisplay()
		display.ProjectName = pMap[v.ProjectCode].Name
		display.MemberAvatar = mMap[v.MemberCode].Avatar
		display.MemberName = mMap[v.MemberCode].Name
		display.TaskName = tMap[v.SourceCode].Name
		list = append(list, display)
	}
	var msgList []*pg.ProjectLogMessage
	copier.Copy(&msgList, list)
	return &pg.ProjectLogResponse{
		List:  msgList,
		Total: total,
	}, nil
}

func (svc *ProjectService) NodeList(context.Context, *pg.ProjectRequest) (*pg.ProjectNodeResponseMessage, error) {
	list, err := svc.nodeDomain.TreeList()
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	var nodes []*pg.ProjectNodeMessage
	copier.Copy(&nodes, list)
	return &pg.ProjectNodeResponseMessage{Nodes: nodes}, nil
}

func (svc *ProjectService) FindProjectByMid(ctx context.Context, msg *pg.ProjectRequest) (*pg.FindProjectByMemberIdResponse, error) {
	isProjectCode := false
	var projectId int64
	if msg.ProjectCode != "" {
		projectId = encrypts.DecryptInt64(msg.ProjectCode)
		isProjectCode = true
	}
	isTaskCode := false
	var taskId int64
	if msg.TaskCode != "" {
		taskId = encrypts.DecryptInt64(msg.TaskCode)
		isTaskCode = true
	}
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if !isProjectCode && isTaskCode {
		projectCode, ok, bError := svc.taskDomain.FindProjectIdByTaskId(taskId)
		if bError != nil {
			return nil, bError
		}
		if !ok {
			return &pg.FindProjectByMemberIdResponse{
				Project:  nil,
				IsOwner:  false,
				IsMember: false,
			}, nil
		}
		projectId = projectCode
		isProjectCode = true
	}
	if isProjectCode {
		//根据projectid和memberid查询
		pm, err := svc.prepo.FindProjectByPidAndMemberId(c, projectId, msg.MemberId)
		if err != nil {
			return nil, model.MySQLError
		}
		if pm == nil {
			return &pg.FindProjectByMemberIdResponse{
				Project:  nil,
				IsOwner:  false,
				IsMember: false,
			}, nil
		}
		projectMessage := &pg.ProjectMessage{}
		copier.Copy(projectMessage, pm)
		isOwner := false
		if pm.IsOwner == 1 {
			isOwner = true
		}
		return &pg.FindProjectByMemberIdResponse{
			Project:  projectMessage,
			IsOwner:  isOwner,
			IsMember: true,
		}, nil
	}
	return &pg.FindProjectByMemberIdResponse{}, nil
}
