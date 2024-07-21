package service

import (
	"context"
	"fmt"
	"time"

	pu "test.com/project-grpc/user_grpc"

	"test.com/project-project/internal/domain"

	"test.com/project-project/config"

	"test.com/common/tms"

	"github.com/jinzhu/copier"

	"go.uber.org/zap"
	"test.com/common/errs"
	"test.com/project-project/internal/model"

	"test.com/common/encrypts"

	pg "test.com/project-grpc/project_grpc"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository"
	"test.com/project-project/internal/repository/dao"
	"test.com/project-project/internal/repository/database"
)

type TaskService struct {
	pg.UnimplementedTaskServiceServer
	cache   repo.Cache
	tran    database.Transaction
	mrepo   *repository.MenuRepository
	prepo   *repository.ProjectRepository
	trepo   *repository.ProjectTemplateRepository
	tsrepo  *repository.TaskStagesRepository
	frepo   *repository.FileRepository
	udomain *domain.UserRpcDomain
	twt     *domain.TaskWorkTimeDomain
}

func NewTaskService() *TaskService {
	return &TaskService{
		cache:   dao.NewRedisCache(),
		tran:    repository.NewTransaction(),
		mrepo:   repository.NewMenuRepository(),
		prepo:   repository.NewProjectRepository(),
		trepo:   repository.NewProjectTemplateRepository(),
		tsrepo:  repository.NewTaskStagesRepository(),
		frepo:   repository.NewFileRepository(),
		udomain: domain.NewUserRpcDomain(),
		twt:     domain.NewTaskWorkTimeDomain(),
	}
}

func (svc *TaskService) TaskStages(ctxs context.Context, msg *pg.TaskRequest) (*pg.TaskStagesResponse, error) {
	pCode := encrypts.DecryptInt64(msg.ProjectCode)
	page := msg.Page
	pageSize := msg.PageSize
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stages, total, err := svc.tsrepo.FindTaskStagesByProjectId(ctx, pCode, page, pageSize)
	if err != nil {
		zap.L().Error("task TaskStages FindTaskStagesByProjectId() failed, error:", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	var tsMessage []*pg.TaskStagesMessage
	copier.Copy(&tsMessage, stages)
	//fmt.Println("tsMessage", tsMessage)
	if tsMessage == nil {
		return &pg.TaskStagesResponse{List: tsMessage, Total: 0}, nil
	}
	stagesMap := model.ToTaskStagesMap(stages)
	for _, v := range tsMessage {
		taksStages := stagesMap[int(v.Id)]
		v.Code = encrypts.EncryptInt64NoErr(int64(v.Id))
		v.CreateTime = tms.FormatByMill(taksStages.CreateTime)
		v.ProjectCode = msg.ProjectCode
	}
	return &pg.TaskStagesResponse{List: tsMessage, Total: total}, nil
}

func (svc *TaskService) MemberProjectList(ctxs context.Context, msg *pg.TaskRequest) (*pg.MemberProjectResponse, error) {
	//1.去ms_project_member表查询用户id
	pCode := encrypts.DecryptInt64(msg.ProjectCode)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	projectMembers, total, err := svc.prepo.FindProjectMemberByPid(ctx, pCode)
	if err != nil {
		zap.L().Error("task MemberProjectList FindProjectMemberByPid() failed, error:", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	//2.根据用户id再到member表查询用户信息
	//调用user grpc服务传入member_id获取用户信息
	if projectMembers == nil || len(projectMembers) <= 0 {
		return &pg.MemberProjectResponse{List: nil, Total: 0}, nil
	}
	var mids []int64

	pmMap := make(map[int64]*model.ProjectMember)
	for _, v := range projectMembers {
		mids = append(mids, v.MemberCode) //v.id := 38
		fmt.Println("v.MemberCode:", v.MemberCode)
		pmMap[v.MemberCode] = v
	}
	//fmt.Println("程序运行到此处，mids:", mids)
	//请求用户信息
	memberList, _, err := svc.udomain.MemberList(mids)
	if err != nil {
		zap.L().Error("task MemberProjectList udomain.MemberList(mids) failed, error:", zap.Error(err))
		return nil, err
	}
	//fmt.Println("程序运行到此处,memberList:", memberList, "pmMap:", pmMap)
	var list []*pg.MemberProjectMessage
	for _, v := range memberList {
		owner := pmMap[v.Mid].IsOwner
		mpm := &pg.MemberProjectMessage{
			MemberCode: v.Mid,
			Name:       v.Name,
			Avatar:     v.Avatar,
			Email:      v.Email,
			Code:       v.Code,
		}
		if v.Mid == owner {
			mpm.IsOwner = config.Owner
		}
		list = append(list, mpm)
	}
	return &pg.MemberProjectResponse{List: list, Total: total}, nil
}

func (svc *TaskService) TaskList(ctxs context.Context, msg *pg.TaskRequest) (*pg.TaskListResponse, error) {
	stageCode := encrypts.DecryptInt64(msg.StageCode)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	taskList, err := svc.tsrepo.FindTaskByStageCode(ctx, int(stageCode))
	if err != nil {
		zap.L().Error("task TaskList FindTaskByStageCode error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	fmt.Println("taskList:member_id:", msg.MemberId)
	var taskDisplayList []*model.TaskDisplay
	var mids []int64
	for _, v := range taskList {
		display := v.ToTaskDisplay()
		if v.Private == 1 {
			tm, err := svc.tsrepo.FindTaskMemberByTaskId(ctx, v.Id, msg.MemberId)
			if err != nil {
				zap.L().Error(" task TaskList FindTaskMemberByTaskId() error", zap.Error(err))
				return nil, errs.GrpcError(model.MySQLError)
			}
			if tm == nil {
				display.CanRead = config.NoCanRead
			} else {
				display.CanRead = config.CanRead
			}
		}
		taskDisplayList = append(taskDisplayList, display)
		mids = append(mids, v.AssignTo)
	}
	//用户信息为空
	if mids == nil || len(mids) <= 0 {
		return &pg.TaskListResponse{List: nil}, nil
	}
	_, memberMap, err := svc.udomain.MemberList(mids)
	if err != nil {
		zap.L().Error(" task TaskList udomain.MemberList(mids) error", zap.Error(err))
		return nil, err
	}
	//for _, v := range memberList {
	//	memberMap[v.Mid] = v
	//}
	for _, v := range taskDisplayList {
		message := memberMap[encrypts.DecryptInt64(v.AssignTo)]
		e := model.Executor{
			Name:   message.Name,
			Avatar: message.Avatar,
		}
		v.Executor = e
	}
	var taskMessageList []*pg.TaskMessage
	copier.Copy(&taskMessageList, taskDisplayList)
	return &pg.TaskListResponse{List: taskMessageList}, nil
}

func (svc *TaskService) SaveTask(ctxs context.Context, msg *pg.TaskRequest) (*pg.TaskMessage, error) {
	//先检查业务
	if msg.Name == "" {
		return nil, errs.GrpcError(model.TaskNameNotNull)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stageCode := encrypts.DecryptInt64(msg.StageCode)
	taskStages, err := svc.tsrepo.FindById(ctx, int(stageCode))
	if err != nil {
		zap.L().Error(" task SaveTask FindById() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if taskStages == nil {
		return nil, errs.GrpcError(model.TaskStagesNotNull)
	}
	projectCode := encrypts.DecryptInt64(msg.ProjectCode)
	project, err := svc.prepo.FindProjectByPId(ctx, projectCode)
	if err != nil {
		zap.L().Error(" task SaveTask FindProjectByPId() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if project.Deleted == config.Deleted || project == nil {
		return nil, errs.GrpcError(model.ProjectAlreadyDeleted)
	}
	maxIdNum, err := svc.tsrepo.FindTaskMaxIdNum(ctx, projectCode)
	if err != nil {
		zap.L().Error(" task SaveTask FindTaskMaxIdNum() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if maxIdNum == nil {
		a := 0
		maxIdNum = &a
	}
	maxSort, err := svc.tsrepo.FindTaskSort(ctx, projectCode, stageCode)
	if err != nil {
		zap.L().Error("task SaveTask FindTaskSort error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if maxSort == nil {
		a := 0
		maxSort = &a
	}
	assignTo := encrypts.DecryptInt64(msg.AssignTo)
	ts := &model.Task{
		Name:        msg.Name,
		CreateTime:  time.Now().UnixMilli(),
		CreateBy:    msg.MemberId,
		AssignTo:    assignTo,
		ProjectCode: projectCode,
		StageCode:   int(stageCode),
		IdNum:       *maxIdNum + 1,
		Private:     project.OpenTaskPrivate,
		Sort:        *maxSort + 65536,
		BeginTime:   time.Now().UnixMilli(),
		EndTime:     time.Now().Add(2 * 24 * time.Hour).UnixMilli(),
	}
	//开启事物，执行插入任务以及任务成员等相关操作
	err = svc.tran.Action(func(conn database.DBConn) error {
		//向任务表插入数据
		err = svc.tsrepo.SaveTask(ctx, conn, ts)
		if err != nil {
			zap.L().Error("task SaveTask SaveTask() error", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		tm := &model.TaskMember{
			MemberCode: assignTo,
			TaskCode:   ts.Id,
			JoinTime:   time.Now().UnixMilli(),
			IsOwner:    config.Owner,
		}
		if assignTo == msg.MemberId {
			tm.IsExecutor = config.Executor
		}
		//向任务成员表插入数据
		err = svc.tsrepo.SaveTaskMember(ctx, conn, tm)
		if err != nil {
			zap.L().Error("task SaveTask SaveTaskMember() error", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	member, err := svc.udomain.MemberInfo(assignTo)
	if err != nil {
		zap.L().Error(" task TaskList udomain.MemberInfo(assignTo) error", zap.Error(err))
		return nil, err
	}
	display := ts.ToTaskDisplay()
	display.Executor = model.Executor{
		Name:   member.Name,
		Avatar: member.Avatar,
		Code:   member.Code,
	}

	//添加任务动态
	err = createProjectLog(svc.prepo, ts.ProjectCode, ts.Id, ts.Name, ts.AssignTo, "create", "task")
	if err != nil {
		zap.L().Error(" task TaskList createProjectLog() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	tm := &pg.TaskMessage{}
	copier.Copy(tm, display)
	return tm, nil
}

func createProjectLog(prepo *repository.ProjectRepository, projectCode int64, taskCode int64, taskName string, toMemberCode int64, logType string, actionType string) error {
	remark := ""
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if logType == "create" {
		remark = "创建了任务"
	}
	pl := &model.ProjectLog{
		MemberCode:  toMemberCode,
		SourceCode:  taskCode,
		Content:     taskName,
		Remark:      remark,
		ProjectCode: projectCode,
		CreateTime:  time.Now().UnixMilli(),
		Type:        logType,
		ActionType:  actionType,
		Icon:        "plus",
		IsComment:   0,
		IsRobot:     0,
	}
	err := prepo.SaveProjectLog(ctx, pl)
	return err
}

func (svc *TaskService) TaskSort(ctxs context.Context, msg *pg.TaskRequest) (*pg.TaskSortResponse, error) {
	preTaskCode := encrypts.DecryptInt64(msg.PreTaskCode)
	toStageCode := encrypts.DecryptInt64(msg.ToStageCode)
	if msg.PreTaskCode == msg.NextTaskCode {
		return &pg.TaskSortResponse{}, nil
	}
	err := svc.sortTask(preTaskCode, msg.NextTaskCode, toStageCode)
	if err != nil {
		return nil, err
	}
	return &pg.TaskSortResponse{}, nil
}

func (svc *TaskService) sortTask(preTaskCode int64, nextTaskCode string, toStageCode int64) error {
	//1. 从小到大排
	//2. 原有的顺序  比如 1 2 3 4 5 4排到2前面去 4的序号在1和2 之间 如果4是最后一个 保证 4比所有的序号都大  如果 排到第一位 直接置为0
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ts, err := svc.tsrepo.FindTaskById(ctx, preTaskCode)
	if err != nil {
		zap.L().Error("task sortTask FindTaskById() error", zap.Error(err))
		return errs.GrpcError(model.MySQLError)
	}
	err = svc.tran.Action(func(conn database.DBConn) error {
		ts.StageCode = int(toStageCode)
		if nextTaskCode != "" {
			//意味着要进行排序的替换
			nextTaskCode := encrypts.DecryptInt64(nextTaskCode)
			next, err := svc.tsrepo.FindTaskById(ctx, nextTaskCode)
			if err != nil {
				zap.L().Error("task sortTask FindTaskById() error", zap.Error(err))
				return errs.GrpcError(model.MySQLError)
			}
			//next.Sort要找到比它小的那个任务
			prepre, err := svc.tsrepo.FindTaskByStageCodeLtSort(ctx, next.StageCode, next.Sort)
			if err != nil {
				zap.L().Error("task sortTask FindTaskByStageCodeLtSort() error", zap.Error(err))
				return errs.GrpcError(model.MySQLError)
			}
			if prepre != nil {
				ts.Sort = (prepre.Sort + next.Sort) / 2
			}
			if prepre == nil {
				ts.Sort = 0
			}
			//sort := ts.Sort
			//ts.Sort = next.Sort
			//next.Sort = sort
			//err = svc.tsrepo.UpdateTaskSort(ctx, conn, next)
			//if err != nil {
			//	zap.L().Error("task TaskSort FindTaskById() error", zap.Error(err))
			//	return errs.GrpcError(model.MySQLError)
			//}
			//isChange = true
		} else {
			maxSort, err := svc.tsrepo.FindTaskSort(ctx, ts.ProjectCode, int64(ts.StageCode))
			if err != nil {
				zap.L().Error("task sortTask FindTaskSort() error", zap.Error(err))
				return errs.GrpcError(model.MySQLError)
			}
			if maxSort == nil {
				a := 0
				maxSort = &a
			}
			ts.Sort = *maxSort + 65536
		}
		//阈值为50
		if ts.Sort < 50 {
			//重置排序
			err = svc.resetSort(toStageCode)
			if err != nil {
				zap.L().Error("task sortTask resetSort() error", zap.Error(err))
				return errs.GrpcError(model.MySQLError)
			}
			return svc.sortTask(preTaskCode, nextTaskCode, toStageCode)
		}
		//if isChange {
		err = svc.tsrepo.UpdateTaskSort(ctx, conn, ts)
		if err != nil {
			zap.L().Error("task TaskSort FindTaskById() error", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		//}
		return nil
	})
	return err
}

func (svc *TaskService) resetSort(code int64) error {
	list, err := svc.tsrepo.FindTaskByStageCode(context.Background(), int(code))
	if err != nil {
		return err
	}
	return svc.tran.Action(func(conn database.DBConn) error {
		iSort := 65536
		for index, v := range list {
			v.Sort = (index + 1) * iSort
			return svc.tsrepo.UpdateTaskSort(context.Background(), conn, v)
		}
		return nil
	})
}

// ReadTask 任务详情
func (svc *TaskService) ReadTask(ctxs context.Context, msg *pg.TaskRequest) (*pg.TaskMessage, error) {
	//1.获取任务id
	//2.查询ms_task表,获取任务信息
	//3.查询项目详情
	//4.查询任务成员详情
	taskCode := encrypts.DecryptInt64(msg.TaskCode)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	taskInfo, err := svc.tsrepo.FindTaskById(ctx, taskCode)
	if err != nil {
		zap.L().Error("task ReadTask FindTaskById() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if taskInfo == nil {
		return &pg.TaskMessage{}, nil
	}
	display := taskInfo.ToTaskDisplay()
	if taskInfo.Private == 1 {
		//代表隐私模式
		taskMember, err := svc.tsrepo.FindTaskMemberByTaskId(ctx, taskInfo.Id, msg.MemberId)
		if err != nil {
			zap.L().Error("task TaskList FindTaskMemberByTaskId() error", zap.Error(err))
			return nil, errs.GrpcError(model.MySQLError)
		}
		if taskMember != nil {
			display.CanRead = config.CanRead
		} else {
			display.CanRead = config.NoCanRead
		}
	}
	pj, err := svc.prepo.FindProjectByPId(ctx, taskInfo.ProjectCode)
	display.ProjectName = pj.Name
	taskStages, err := svc.tsrepo.FindById(ctx, taskInfo.StageCode)
	display.StageName = taskStages.Name
	// in ()
	//memberMessage, err := rpc.UserServiceClient.FindMemberByMemId(ctx, &pu.UserRequest{MemberId: taskInfo.AssignTo})
	member, err := svc.udomain.MemberInfo(taskInfo.AssignTo)
	if err != nil {
		zap.L().Error("task TaskList udomain.MemberInfo(taskInfo.AssignTo) error", zap.Error(err))
		return nil, err
	}
	e := model.Executor{
		Name:   member.Name,
		Avatar: member.Avatar,
	}
	display.Executor = e
	var taskMessage = &pg.TaskMessage{}
	copier.Copy(taskMessage, display)
	return taskMessage, nil
}

func (svc *TaskService) MyTaskList(ctxs context.Context, msg *pg.TaskRequest) (*pg.MyTaskListResponse, error) {
	var tsList []*model.Task
	var err error
	var total int64
	page := msg.Page
	pageSize := msg.PageSize
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	switch msg.TaskType {
	case 1:
		//我执行的
		tsList, total, err = svc.tsrepo.FindTaskByAssignTo(ctx, msg.MemberId, int(msg.Type), page, pageSize)
		if err != nil {
			zap.L().Error("task MyTaskList FindTaskByAssignTo() failed", zap.Error(err))
			return nil, errs.GrpcError(model.MySQLError)
		}
	case 2:
		//我参加的
		tsList, total, err = svc.tsrepo.FindTaskByMemberCode(ctx, msg.MemberId, int(msg.Type), page, pageSize)
		if err != nil {
			zap.L().Error("task MyTaskList FindTaskByMemberCode() failed", zap.Error(err))
			return nil, errs.GrpcError(model.MySQLError)
		}
	case 3:
		//我创建的
		tsList, total, err = svc.tsrepo.FindTaskByCreateBy(ctx, msg.MemberId, int(msg.Type), page, pageSize)
		if err != nil {
			zap.L().Error("task MyTaskList FindTaskByCreateBy() failed", zap.Error(err))
			return nil, errs.GrpcError(model.MySQLError)
		}
	default:
		return nil, errs.GrpcError(model.ErrorServerBusy)
	}
	if tsList == nil || len(tsList) <= 0 {
		return &pg.MyTaskListResponse{List: nil, Total: 0}, nil
	}
	var pids []int64
	var mids []int64
	for _, v := range tsList {
		pids = append(pids, v.ProjectCode)
		mids = append(mids, v.AssignTo)
	}
	pListChan := make(chan []*model.Project)
	defer close(pListChan)
	mListChan := make(chan []*pu.MemberResponse)
	defer close(mListChan)
	go func() {
		pList, err := svc.prepo.FindProjectByPIds(ctx, pids)
		if err != nil {
			zap.L().Error("task MyTaskList FindProjectByIds() failed", zap.Error(err))
		}
		pListChan <- pList
	}()
	go func() {
		//mList, err := rpc.UserServiceClient.FindMemberInfoByMIds(ctxs, &pu.UserRequest{
		//	MIds: mids,
		//})
		mList, _, err := svc.udomain.MemberList(mids)
		if err != nil {
			zap.L().Error("task MyTaskList udomain.MemberList(mids) failed", zap.Error(err))
		}
		mListChan <- mList
	}()
	pList := <-pListChan
	projectMap := model.ToProjectMap(pList)
	mList := <-mListChan
	mMap := make(map[int64]*pu.MemberResponse)
	for _, v := range mList {
		mMap[v.Mid] = v
	}
	var mtdList []*model.MyTaskDisplay
	for _, v := range tsList {
		memberMessage := mMap[v.AssignTo]
		name := memberMessage.Name
		avatar := memberMessage.Avatar
		mtd := v.ToMyTaskDisplay(projectMap[v.ProjectCode], name, avatar)
		mtdList = append(mtdList, mtd)
	}
	var myMsgs []*pg.MyTaskMessage
	copier.Copy(&myMsgs, mtdList)
	return &pg.MyTaskListResponse{List: myMsgs, Total: total}, nil
}

func (svc *TaskService) ListTaskMember(ctxs context.Context, msg *pg.TaskRequest) (*pg.TaskMemberList, error) {
	taskCode := encrypts.DecryptInt64(msg.TaskCode)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	taskMemberPage, total, err := svc.tsrepo.FindTaskMemberByTaskIdAndPage(ctx, taskCode, msg.Page, msg.PageSize)
	if err != nil {
		zap.L().Error("task TaskList FindTaskMemberPage() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	var mids []int64
	for _, v := range taskMemberPage {
		mids = append(mids, v.MemberCode)
	}

	_, mMap, err := svc.udomain.MemberList(mids)
	if err != nil {
		zap.L().Error("task TaskList udomain.MemberList() error", zap.Error(err))
		return nil, err
	}
	var taskMemeberMemssages []*pg.TaskMemberMessage
	for _, v := range taskMemberPage {
		tm := &pg.TaskMemberMessage{}
		tm.Code = encrypts.EncryptInt64NoErr(v.MemberCode)
		tm.Id = v.Id
		message := mMap[v.MemberCode]
		tm.Name = message.Name
		tm.Avatar = message.Avatar
		tm.IsExecutor = int32(v.IsExecutor)
		tm.IsOwner = int32(v.IsOwner)
		taskMemeberMemssages = append(taskMemeberMemssages, tm)
	}
	return &pg.TaskMemberList{List: taskMemeberMemssages, Total: total}, nil
}

func (svc *TaskService) TaskLog(ctxs context.Context, msg *pg.TaskRequest) (*pg.TaskLogList, error) {
	taskCode := encrypts.DecryptInt64(msg.TaskCode)
	all := msg.All
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var list []*model.ProjectLog
	var total int64
	var err error
	if all == 1 {
		//显示全部
		list, total, err = svc.prepo.FindProjectLogByTaskCode(c, taskCode, int(msg.Comment))
	}
	if all == 0 {
		//分页
		list, total, err = svc.prepo.FindProjectLogByTaskCodePage(c, taskCode, int(msg.Comment), int(msg.Page), int(msg.PageSize))
	}
	if err != nil {
		zap.L().Error("task TaskLog FindLogByTaskCodePage() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if total == 0 {
		return &pg.TaskLogList{}, nil
	}
	var displayList []*model.ProjectLogDisplay
	var mIdList []int64
	for _, v := range list {
		mIdList = append(mIdList, v.MemberCode)
	}
	//messageList, err := rpc.UserServiceClient.FindMemberInfoByMIds(c, &pu.UserRequest{MIds: mIdList})
	//mMap := make(map[int64]*pu.MemberResponse)
	//for _, v := range messageList.List {
	//	mMap[v.Mid] = v
	//}
	_, mMap, err := svc.udomain.MemberList(mIdList)
	if err != nil {
		zap.L().Error("task TaskLog udomain.MemberList() error", zap.Error(err))
		return nil, err
	}
	for _, v := range list {
		display := v.ToDisplay()
		message := mMap[v.MemberCode]
		m := model.Member{}
		m.Name = message.Name
		m.Id = message.Id
		m.Avatar = message.Avatar
		m.Code = message.Code
		display.Member = m
		displayList = append(displayList, display)
	}
	var l []*pg.TaskLog
	copier.Copy(&l, displayList)
	return &pg.TaskLogList{List: l, Total: total}, nil
}

func (svc *TaskService) TaskWorkTimeList(ctx context.Context, msg *pg.TaskRequest) (*pg.TaskWorkTimeResponse, error) {
	taskCode := encrypts.DecryptInt64(msg.TaskCode)
	list, err := svc.twt.TaskWorkTimeList(taskCode)
	if err != nil {
		zap.L().Error("task TaskWorkTimeList twt.TaskWorkTimeList(taskCode) error", zap.Error(err))
		return nil, err
	}
	var l []*pg.TaskWorkTime
	copier.Copy(&l, list)
	return &pg.TaskWorkTimeResponse{List: l, Total: int64(len(l))}, nil
}

func (svc *TaskService) SaveTaskWorkTime(ctx context.Context, msg *pg.TaskRequest) (*pg.SaveTaskWorkTimeResponse, error) {
	tmt := &model.TaskWorkTime{}
	tmt.BeginTime = msg.BeginTime
	tmt.Num = int(msg.Num)
	tmt.Content = msg.Content
	tmt.TaskCode = encrypts.DecryptInt64(msg.TaskCode)
	tmt.MemberCode = msg.MemberId
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := svc.tsrepo.SaveTaskWorkTime(c, tmt)
	if err != nil {
		zap.L().Error("task SaveTaskWorkTime SaveTaskWorkTime() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	return &pg.SaveTaskWorkTimeResponse{}, nil
}

func (svc *TaskService) SaveTaskFile(ctxs context.Context, msg *pg.TaskFileReqMessage) (*pg.TaskFileResponse, error) {
	taskCode := encrypts.DecryptInt64(msg.TaskCode)
	fmt.Println("SaveTaskFile taskCode:", taskCode)
	//存file表
	f := &model.File{
		PathName:         msg.PathName,
		Title:            msg.FileName,
		Extension:        msg.Extension,
		Size:             int(msg.Size),
		ObjectType:       "",
		OrganizationCode: encrypts.DecryptInt64(msg.OrganizationCode),
		TaskCode:         taskCode,
		ProjectCode:      encrypts.DecryptInt64(msg.ProjectCode),
		CreateBy:         msg.MemberId,
		CreateTime:       time.Now().UnixMilli(),
		Downloads:        0,
		Extra:            "",
		Deleted:          config.NoDeleted,
		FileType:         msg.FileType,
		FileUrl:          msg.FileUrl,
		DeletedTime:      0,
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	err := svc.tran.Action(func(conn database.DBConn) error {
		err := svc.frepo.SaveFileInfo(ctx, conn, f)
		if err != nil {
			zap.L().Error("task SaveTaskFile SaveFileInfo() error", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		//存入source_link
		sl := &model.SourceLink{
			SourceType:       "file",
			SourceCode:       f.Id,
			LinkType:         "task",
			LinkCode:         taskCode,
			OrganizationCode: encrypts.DecryptInt64(msg.OrganizationCode),
			CreateBy:         msg.MemberId,
			CreateTime:       time.Now().UnixMilli(),
			Sort:             0,
		}
		err = svc.frepo.SaveSourceLinkInfo(ctx, conn, sl)
		if err != nil {
			zap.L().Error("task SaveTaskFile SaveFileInfo() error", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pg.TaskFileResponse{}, nil
}

func (svc *TaskService) TaskSources(ctxs context.Context, msg *pg.TaskRequest) (*pg.TaskSourceResponse, error) {
	taskCode := encrypts.DecryptInt64(msg.TaskCode)
	fmt.Println("TaskSources() taskCode:", taskCode, "msg.TaskCode:", msg.TaskCode)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()
	sourceLinks, err := svc.frepo.FindSourceLinksByTaskCode(ctx, taskCode)
	if err != nil {
		zap.L().Error("task TaskSources FindSourcesByTaskCode() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if len(sourceLinks) == 0 {
		return &pg.TaskSourceResponse{}, nil
	}
	var fIdList []int64
	for _, v := range sourceLinks {
		fIdList = append(fIdList, v.SourceCode)
	}
	files, err := svc.frepo.FindFileByIds(ctx, fIdList)
	if err != nil {
		zap.L().Error("task TaskSources FindFileByIds() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	fMap := make(map[int64]*model.File)
	for _, v := range files {
		fMap[v.Id] = v
	}
	var list []*model.SourceLinkDisplay
	for _, v := range sourceLinks {
		list = append(list, v.ToDisplay(fMap[v.SourceCode]))
	}
	var slMsg []*pg.TaskSourceMessage
	copier.Copy(&slMsg, list)
	return &pg.TaskSourceResponse{List: slMsg}, nil
}

func (svc *TaskService) CreateTaskComment(ctxs context.Context, msg *pg.TaskRequest) (*pg.CreateCommentResponse, error) {
	taskCode := encrypts.DecryptInt64(msg.TaskCode)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	taskById, err := svc.tsrepo.FindTaskById(ctx, taskCode)
	if err != nil {
		zap.L().Error("task CreateTaskComment FindFileByIds() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	pl := &model.ProjectLog{
		MemberCode:   msg.MemberId,
		Content:      msg.CommentContent,
		Remark:       msg.CommentContent,
		Type:         "createComment",
		CreateTime:   time.Now().UnixMilli(),
		SourceCode:   taskCode,
		ActionType:   "task",
		ToMemberCode: 0,
		IsComment:    config.Comment,
		ProjectCode:  taskById.ProjectCode,
		Icon:         "plus",
		IsRobot:      0,
	}
	err = svc.prepo.SaveProjectLog(ctx, pl)
	if err != nil {
		zap.L().Error("task CreateTaskComment SaveProjectLog() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	return &pg.CreateCommentResponse{}, nil
}

func (svc *TaskService) CreateTaskList(ctxs context.Context, msg *pg.TaskRequest) (*pg.TaskStagesMessage, error) {
	pCode := encrypts.DecryptInt64(msg.ProjectCode)
	fmt.Println("pCode:", pCode)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	//查出任务列表最后一个sort
	ts, err := svc.tsrepo.FindTaskStagesByPcode(ctx, pCode)
	if err != nil {
		zap.L().Error("task CreateTaskList FindTaskStagesByPcode() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	//向ms_task_stages表插入数据
	taskStage := &model.TaskStages{
		ProjectCode: pCode,
		Name:        msg.Name,
		Sort:        ts.Sort + 1,
		Description: "",
		CreateTime:  time.Now().UnixMilli(),
		Deleted:     config.NoDeleted,
	}
	err = svc.tsrepo.CreateTaskStageByStrcut(ctx, taskStage)
	if err != nil {
		zap.L().Error("task CreateTaskList CreateTaskStage() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	return &pg.TaskStagesMessage{
		Code:        encrypts.EncryptInt64NoErr(int64(taskStage.Id)),
		Name:        taskStage.Name,
		Description: taskStage.Description,
		Sort:        int32(taskStage.Sort),
		ProjectCode: msg.ProjectCode,
		CreateTime:  tms.FormatByMill(taskStage.CreateTime),
		Deleted:     int32(taskStage.Deleted),
	}, nil
}

func (svc *TaskService) DeleteTaskStages(ctxs context.Context, msg *pg.TaskRequest) (*pg.UpdateTaskStagesResponse, error) {
	taskStageCode := encrypts.DecryptInt64(msg.ToStageCode)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	//根据任务列表id删除任务列表
	err := svc.tsrepo.DeleteTaskStagesByCode(ctx, taskStageCode)
	if err != nil {
		zap.L().Error("task DeleteTaskStages DeleteTaskStagesByCode() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	return &pg.UpdateTaskStagesResponse{}, nil
}

func (svc *TaskService) EditTaskStages(ctxs context.Context, msg *pg.TaskRequest) (*pg.UpdateTaskStagesResponse, error) {
	taskStageCode := encrypts.DecryptInt64(msg.ToStageCode)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ts := &model.TaskStages{
		Name: msg.Name,
	}
	err := svc.tsrepo.EditTaskStagesByStruct(ctx, taskStageCode, ts)
	if err != nil {
		zap.L().Error("task DeleteTaskStages DeleteTaskStagesByCode() error", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	return &pg.UpdateTaskStagesResponse{}, nil
}
