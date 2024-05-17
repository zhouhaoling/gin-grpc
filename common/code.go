package common

// res_code 业务码

type ResCode int64

const (
	CodeSuccess       ResCode = 2000 + iota //成功
	CodeNoLegalMobile                       //手机不合法
	CodeServerBusy                          //服务繁忙
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:       "success",
	CodeNoLegalMobile: "手机不合法",
	CodeServerBusy:    "服务繁忙",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
