package common

// res_code 业务码

type ResCode int64

const (
	CodeSuccess       ResCode = 200         //成功
	CodeNoLegalMobile ResCode = 2001 + iota //手机不合法
	CodeServerBusy                          //服务繁忙
	CodeInvalidParams
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:       "success",
	CodeNoLegalMobile: "手机不合法",
	CodeServerBusy:    "服务繁忙",
	CodeInvalidParams: "请求参数错误",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
