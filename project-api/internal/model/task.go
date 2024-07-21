package model

type TaskStagesResp struct {
	Name         string `json:"name"`
	ProjectCode  string `json:"project_code"`
	Sort         int    `json:"sort"`
	Description  string `json:"description"`
	CreateTime   string `json:"create_time"`
	Code         string `json:"code"`
	Deleted      int    `json:"deleted"`
	TasksLoading bool   `json:"tasksLoading"`
	FixedCreator bool   `json:"fixedCreator"`
	ShowTaskCard bool   `json:"showTaskCard"`
	Tasks        []int  `json:"tasks"`
	DoneTasks    []int  `json:"doneTasks"`
	UnDoneTasks  []int  `json:"unDoneTasks"`
}

type TaskDisplay struct {
	ProjectCode   string   `json:"project_code"`
	Name          string   `json:"name"`
	Pri           int      `json:"pri"`
	ExecuteStatus string   `json:"execute_status"`
	Description   string   `json:"description"`
	CreateBy      string   `json:"create_by"`
	DoneBy        string   `json:"done_by"`
	DoneTime      string   `json:"done_time"`
	CreateTime    string   `json:"create_time"`
	AssignTo      string   `json:"assign_to"`
	Deleted       int      `json:"deleted"`
	StageCode     string   `json:"stage_code"`
	TaskTag       string   `json:"task_tag"`
	Done          int      `json:"done"`
	BeginTime     string   `json:"begin_time"`
	EndTime       string   `json:"end_time"`
	RemindTime    string   `json:"remind_time"`
	Pcode         string   `json:"pcode"`
	Sort          int      `json:"sort"`
	Like          int      `json:"like"`
	Star          int      `json:"star"`
	DeletedTime   string   `json:"deleted_time"`
	Private       int      `json:"private"`
	IdNum         int      `json:"id_num"`
	Path          string   `json:"path"`
	Schedule      int      `json:"schedule"`
	VersionCode   string   `json:"version_code"`
	FeaturesCode  string   `json:"features_code"`
	WorkTime      int      `json:"work_time"`
	Status        int      `json:"status"`
	Code          string   `json:"code"`
	CanRead       int      `json:"canRead"`
	HasUnDone     int      `json:"hasUnDone"`
	ParentDone    int      `json:"parentDone"`
	HasComment    int      `json:"hasComment"`
	HasSource     int      `json:"hasSource"`
	Executor      Executor `json:"executor"`
	PriText       string   `json:"priText"`
	StatusText    string   `json:"statusText"`
	Liked         int      `json:"liked"`
	Stared        int      `json:"stared"`
	Tags          []int    `json:"tags"`
	ChildCount    []int    `json:"childCount"`
	ProjectName   string   `json:"projectName"`
	StageName     string   `json:"stageName"`
}

type Executor struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
type TaskSaveReq struct {
	Name        string `form:"name"`
	StageCode   string `form:"stage_code"`
	ProjectCode string `form:"project_code"`
	AssignTo    string `form:"assign_to"`
}

type TaskSortReq struct {
	PreTaskCode  string `form:"preTaskCode"`
	NextTaskCode string `form:"nextTaskCode"`
	ToStageCode  string `form:"toStageCode"`
}

type MyTaskReq struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"pageSize"`
	TaskType int   `form:"taskType" json:"task_type"` //任务类型
	Type     int   `form:"type" json:"type"`          //0未完成 1已完成
}

type MyTaskDisplay struct {
	ProjectCode        string      `json:"project_code"`
	Name               string      `json:"name"`
	Pri                int         `json:"pri"`
	ExecuteStatus      string      `json:"execute_status"`
	Description        string      `json:"description"`
	CreateBy           string      `json:"create_by"`
	DoneBy             string      `json:"done_by"`
	DoneTime           string      `json:"done_time"`
	CreateTime         string      `json:"create_time"`
	AssignTo           string      `json:"assign_to"`
	Deleted            int         `json:"deleted"`
	StageCode          string      `json:"stage_code"`
	TaskTag            string      `json:"task_tag"`
	Done               int         `json:"done"`
	BeginTime          string      `json:"begin_time"`
	EndTime            string      `json:"end_time"`
	RemindTime         string      `json:"remind_time"`
	Pcode              string      `json:"pcode"`
	Sort               int         `json:"sort"`
	Like               int         `json:"like"`
	Star               int         `json:"star"`
	DeletedTime        string      `json:"deleted_time"`
	Private            int         `json:"private"`
	IdNum              int         `json:"id_num"`
	Path               string      `json:"path"`
	Schedule           int         `json:"schedule"`
	VersionCode        string      `json:"version_code"`
	FeaturesCode       string      `json:"features_code"`
	WorkTime           int         `json:"work_time"`
	Status             int         `json:"status"`
	Code               string      `json:"code"`
	ProjectName        string      `json:"project_name"`
	Cover              string      `json:"cover"`
	AccessControlType  string      `json:"access_control_type"`
	WhiteList          string      `json:"white_list"`
	Order              int         `json:"order"`
	TemplateCode       string      `json:"template_code"`
	OrganizationCode   string      `json:"organization_code"`
	Prefix             string      `json:"prefix"`
	OpenPrefix         int         `json:"open_prefix"`
	Archive            int         `json:"archive"`
	ArchiveTime        string      `json:"archive_time"`
	OpenBeginTime      int         `json:"open_begin_time"`
	OpenTaskPrivate    int         `json:"open_task_private"`
	TaskBoardTheme     string      `json:"task_board_theme"`
	AutoUpdateSchedule int         `json:"auto_update_schedule"`
	HasUnDone          int         `json:"hasUnDone"`
	ParentDone         int         `json:"parentDone"`
	PriText            string      `json:"priText"`
	Executor           Executor    `json:"executor"`
	ProjectInfo        ProjectInfo `json:"projectInfo"`
}

type ProjectInfo struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type TaskMember struct {
	Id                int64  `json:"id"`
	Name              string `json:"name"`
	Avatar            string `json:"avatar"`
	Code              int64  `json:"code"`
	IsExecutor        int    `json:"is_executor"`
	IsOwner           int    `json:"is_owner"`
	MemberAccountCode string `json:"member_account_code"`
}

type TaskWorkTime struct {
	Id         int64  `json:"id"`
	TaskCode   string `json:"task_code"`
	MemberCode string `json:"member_code"`
	CreateTime string `json:"create_time"`
	Content    string `json:"content"`
	BeginTime  string `json:"begin_time"`
	Num        int    `json:"num"`
	Code       string `json:"code"`
	Member     Member `json:"member"`
}

type SaveTaskWorkTimeReq struct {
	TaskCode  string `json:"task_code" form:"taskCode"`
	Content   string `form:"content"`
	Num       int    `form:"num"`
	BeginTime string `form:"beginTime"`
}

type CommentReq struct {
	TaskCode string   `form:"taskCode"`
	Comment  string   `form:"comment"`
	Mentions []string `form:"mentions"`
}
