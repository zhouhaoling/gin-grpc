package model

import (
	"test.com/common"
	"test.com/common/errs"
)

var (
	RedisError      = errs.NewError(999, "redis错误")
	MySQLError      = errs.NewError(998, "数据库错误")
	ErrorServerBusy = errs.NewError(common.CodeServerBusy, "服务繁忙")
	// 1 开头 代表用户的业务码
	NoLegalMobile      = errs.NewError(12001, "手机不合法")
	ErrorCaptcha       = errs.NewError(12002, "验证码不合法")
	CaptchaNotExist    = errs.NewError(12003, "验证码已过期")
	EmailExist         = errs.NewError(12004, "邮箱已经存在")
	MobileExist        = errs.NewError(12005, "手机号已经存在")
	MobileNotExist     = errs.NewError(12006, "手机号不存在")
	AccountAndPwdError = errs.NewError(12007, "账号或密码错误")
	OrgNotExist        = errs.NewError(12008, "组织信息不存在")
)
