package model

import (
	"test.com/common/errs"
)

var (
	NoLegalMobile = errs.NewError(2001, "手机不合法")
)
