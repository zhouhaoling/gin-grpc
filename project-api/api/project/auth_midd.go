package project

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"test.com/common"
	"test.com/common/errs"
)

var ignores = []string{
	"project/login/register",
	"project/login",
	"project/login/getCaptcha",
	"project/organization",
	"project/auth/apply"}

func Auth() func(*gin.Context) {
	return func(c *gin.Context) {
		zap.L().Info("开始做授权认证")
		uri := c.Request.RequestURI
		for _, v := range ignores {
			if strings.Contains(uri, v) {
				c.Next()
				return
			}
		}
		//判断此uri是否在用户的授权列表中
		res := common.NewResponseData()
		a := NewHandlerAuth()
		nodes, err := a.GetAuthNodes(c)
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			res.ResponseErrorWithMsg(c, code, msg)
			c.Abort()
			return
		}
		for _, v := range nodes {
			if strings.Contains(uri, v) {
				c.Next()
				return
			}
		}
		//TODO 未补充，判断是否在权限控制的路径列表中，如果在就继续，不在就返回
		res.ResponseErrorWithMsg(c, common.CodeNotAuth, "没有权限访问")
		return
	}
}
