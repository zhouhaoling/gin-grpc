package errs

import (
	"fmt"
	"test.com/common"
)

type BError struct {
	Code common.ResCode
	Msg  string
}

func (e *BError) Error() string {
	return fmt.Sprintf("code:%v,msg:%s", e.Code, e.Msg)
}

func NewError(code common.ResCode, msg string) *BError {
	return &BError{
		Code: code,
		Msg:  msg,
	}
}
