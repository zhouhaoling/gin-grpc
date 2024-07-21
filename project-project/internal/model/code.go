package model

import (
	"test.com/common"
	"test.com/common/errs"
)

var (
	RedisError      = errs.NewError(999, "redis错误")
	MySQLError      = errs.NewError(998, "数据库错误")
	ParamsError     = errs.NewError(401, "参数错误")
	ErrorServerBusy = errs.NewError(common.CodeServerBusy, "服务繁忙")
	// 2 开头 代表项目的业务码
	TaskNameNotNull       = errs.NewError(22001, "任务名称不能为空")
	TaskStagesNotNull     = errs.NewError(22002, "任务步骤不存在")
	ProjectAlreadyDeleted = errs.NewError(22003, "项目已经删除")
)
