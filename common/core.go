package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponseData() *ResponseData {
	return &ResponseData{}
}

func (r *ResponseData) Response(data any) *ResponseData {
	r.Code = CodeSuccess
	r.Msg = CodeSuccess.Msg()
	r.Data = data
	return r
}

// ResponseSuccess 返回成功的响应
func (r *ResponseData) ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

// ResponseSuccessWithMsg 自定义返回成功的响应
func (r *ResponseData) ResponseSuccessWithMsg(c *gin.Context, msg, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

// ResponseError 返回失败的响应
func (r *ResponseData) ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// ResponseErrorWithMsg 自定义返回失败的响应
func (r *ResponseData) ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
