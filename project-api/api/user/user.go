package user

import (
	"context"
	"time"

	"test.com/common/errs"

	"github.com/gin-gonic/gin"
	"test.com/common"
	ug "test.com/project-user/user_grpc"
)

type HandlerUser struct {
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}

func (u *HandlerUser) getCaptcha(c *gin.Context) {
	res := common.NewResponseData()
	mobile := c.PostForm("mobile")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := LoginServiceClient.GetCaptcha(ctx, &ug.CaptchaRequest{Mobile: mobile})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	//返回成功的响应码
	res.ResponseSuccess(c, response.Code)
}
