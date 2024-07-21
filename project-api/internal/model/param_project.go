package model

import "github.com/gin-gonic/gin"

type ParamProject struct {
	Name         string `json:"name" form:"name"`
	TemplateCode string `json:"templateCode" form:"templateCode"`
	Description  string `json:"description" form:"description"`
	Id           int    `json:"id" form:"id"`
}

func (p *ParamProject) Bind(c *gin.Context) error {
	err := c.ShouldBind(p)
	if err != nil {
		return err
	}
	return nil
}

type ParamProjectEdit struct {
	ProjectCode        string `json:"projectCode" form:"projectCode"`                   //项目id
	Name               string `json:"name" form:"name"`                                 //姓名
	Description        string `json:"description" form:"description"`                   //项目描述
	Cover              string `json:"cover" form:"cover"`                               //项目封面
	Private            int    `json:"private" form:"private"`                           //1私有 0公有
	Prefix             string `json:"prefix" form:"prefix"`                             //项目前缀
	TaskBoardTheme     string `json:"task_board_theme" form:"task_board_theme"`         //看板风格
	OpenPrefix         int    `json:"open_prefix" form:"open_prefix"`                   //是否开启项目前缀 默认0 不开启
	OpenBeginTime      int    `json:"open_begin_time" form:"open_begin_time"`           //是否开启项目开始时间 默认0 不开启 1 开启
	OpenTaskPrivate    int    `json:"open_task_private" form:"open_task_private"`       //是否开启新任务隐私模式 0 不开启 1 开启
	Schedule           int    `json:"schedule" form:"schedule"`                         //进度
	AutoUpdateSchedule int    `json:"auto_update_schedule" form:"auto_update_schedule"` //自动更新进度 0 不开启 1 开启
}

func (p *ParamProjectEdit) Bind(c *gin.Context) error {
	return c.ShouldBind(p)
}

type TaskStagesSaveReq struct {
	ProjectCode string `json:"projectCode" form:"projectCode"`
	Name        string `json:"name" form:"name"`
}
