package model

import (
	"test.com/common"
	"test.com/common/errs"
)

var (
	RedisError      = errs.NewError(999, "redis错误")
	MySQLError      = errs.NewError(998, "数据库错误")
	ErrorServerBusy = errs.NewError(common.CodeServerBusy, "服务繁忙")
	// 2 开头 代表项目的业务码

)
