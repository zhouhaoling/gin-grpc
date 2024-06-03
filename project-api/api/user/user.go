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

func (u *HandlerUser) userRegister(c *gin.Context) {
	//1.将参数绑定到结构体中
	//2.校验参数
	//3.调用grpc服务的用户注册服务
	//5.grpc服务中：判断账号是否存在，邮箱是否重复，手机号是否重复
	//6.生产uuid(雪花算法），将密码加密，将数据插入到数据库中
	//7.返回操作后的响应码和err
	//8.根据grpc返回的err来进行判断，返回给前端什么样的响应
}

func (u *HandlerUser) userLogin(c *gin.Context) {
	//1.将参数绑定到结构体中
	//2.校验参数，比如验证码是否为6位数字，密码是否为字母和数字组成等
	//3.调用grpc服务的用户登录服务
	//4.grpc服务中，先查询用户是否存在
	//5.判断验证码是否正确,判断验证码是否正确
	//6.获取用户信息.
}
