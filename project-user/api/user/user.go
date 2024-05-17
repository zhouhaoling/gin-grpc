package user

import (
	"log"
	"net/http"
	"time"

	"test.com/project-user/tools"

	"test.com/project-user/pkg/model"

	"github.com/gin-gonic/gin"
	"test.com/common"
)

type HandlerUser struct {
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}

func (u *HandlerUser) getCaptcha(ctx *gin.Context) {
	rsp := common.NewResult()
	//1.获取参数
	mobile := ctx.PostForm("mobile")
	//2.校验参数
	if !common.VerifyModel(mobile) {
		ctx.JSON(http.StatusOK, rsp.Fail(model.NoLegalMobile, "手机号格式不正确"))
		return
	}
	//3.生成验证码
	code := tools.GetVerifyCode()
	//4.调用短信验证平台
	go func() {
		time.Sleep(2 * time.Second)
		log.Println("调用短信验证平台成功,发送短信")
		//存储验证码到redis中,并设置过期时间
		log.Printf("存入成功：register_%s:%s", mobile, code)
	}()
	ctx.JSON(http.StatusOK, rsp.Success(code))
}
