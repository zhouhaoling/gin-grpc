package model

type TaskLogReq struct {
	TaskCode string `form:"taskCode"`
	PageSize int    `form:"pageSize"`
	Page     int    `form:"page"`
	All      int    `form:"all"`
	Comment  int    `form:"comment"`
}
type ProjectLogDisplay struct {
	Id           int64   `json:"id"`
	MemberCode   string  `json:"member_code"`
	Content      string  `json:"content"`
	Remark       string  `json:"remark"`
	Type         string  `json:"type"`
	CreateTime   string  `json:"create_time"`
	SourceCode   string  `json:"source_code"`
	ActionType   string  `json:"action_type"`
	ToMemberCode string  `json:"to_member_code"`
	IsComment    int     `json:"is_comment"`
	ProjectCode  string  `json:"project_code"`
	Icon         string  `json:"icon"`
	IsRobot      int     `json:"is_robot"`
	Member       MemberL `json:"member"`
}

type MemberL struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Avatar string `json:"avatar"`
}

type ProjectLog struct {
	Content      string `json:"content"`
	Remark       string `json:"remark"`
	CreateTime   string `json:"create_time"`
	SourceCode   string `json:"source_code"`
	IsComment    int    `json:"is_comment"`
	ProjectCode  string `json:"project_code"`
	ProjectName  string `json:"project_name"`
	MemberAvatar string `json:"member_avatar"`
	MemberName   string `json:"member_name"`
	TaskName     string `json:"task_name"`
}