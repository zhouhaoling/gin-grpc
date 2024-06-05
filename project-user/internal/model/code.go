package model

import (
	"test.com/common/errs"
)

var (
	RedisError    = errs.NewError(999, "redis错误")
	MySQLError    = errs.NewError(998, "数据库错误")
	NoLegalMobile = errs.NewError(2001, "手机不合法")
	ErrorCaptcha  = errs.NewError(2002, "验证码不合法")
)
