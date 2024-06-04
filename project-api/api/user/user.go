package user

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"test.com/project-api/internal/model"

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
	//3.调用user grpc服务 获取响应
	//4.返回结果
	res := common.NewResponseData()
	var param model.ParamRegister
	if err := c.ShouldBindJSON(&param); err != nil {
		zap.L().Error("register with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			res.ResponseError(c, common.CodeInvalidParams)
			return
		}
		res.ResponseErrorWithMsg(c, common.CodeInvalidParams, errs.Translate(common.Trans))
		return
	}
	//校验参数
	if err := param.Verify(); err != nil {
		res.ResponseErrorWithMsg(c, common.CodeInvalidParams, err.Error())
		return
	}
	//调用grpc服务
	err := UserServiceClient.Register()
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	//返回成功的响应码
	res.ResponseSuccess(c, common.CodeSuccess)
}

func (u *HandlerUser) userLogin(c *gin.Context) {
	//1.将参数绑定到结构体中
	//2.校验参数，比如验证码是否为6位数字，密码是否为字母和数字组成等
	//3.调用grpc服务的用户登录服务
	//4.grpc服务中，先查询用户是否存在
	//5.判断验证码是否正确,判断验证码是否正确
	//6.获取用户信息.
}
