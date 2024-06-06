package model

import (
	"test.com/common/errs"
)

var (
	RedisError = errs.NewError(999, "redis错误")
	MySQLError = errs.NewError(998, "数据库错误")
	// 1 开头 代表用户的业务码
	NoLegalMobile = errs.NewError(12001, "手机不合法")
	ErrorCaptcha  = errs.NewError(12002, "验证码不合法")
	EmailExist    = errs.NewError(12003, "邮箱已经存在")
	MobileExist   = errs.NewError(12004, "手机已经存在")
)
