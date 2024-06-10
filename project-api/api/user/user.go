package user

import (
	"context"
	"fmt"
	"time"

	ug "test.com/project-grpc/user_grpc"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"test.com/project-api/internal/model"

	"test.com/common/errs"

	"github.com/gin-gonic/gin"
	"test.com/common"
)

type HandlerUser struct {
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}

func (u *HandlerUser) getCaptcha(c *gin.Context) {
	res := common.NewResponseData()
	mobile := c.PostForm("mobile")
	fmt.Println("mobile:", mobile)
	if flag := common.VerifyModel(mobile); flag != true {
		res.ResponseErrorWithMsg(c, common.CodeInvalidParams, "手机号格式不正确")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	response, err := loginServiceClient.GetCaptcha(ctx, &ug.CaptchaRequest{Mobile: mobile})
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
	if err := c.ShouldBind(&param); err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//两种方法，将param的值赋值给msg,使用第一种要求结构体的字段名都一致
	//msg := &ug.RegisterRequest{}
	//if err := copier.Copy(msg, param); err != nil {
	//	res.ResponseErrorWithMsg(c, common.CodeServerBusy, "copy有误")
	//	return
	//}
	msg := &ug.RegisterRequest{
		Mobile:   param.Mobile,
		Password: param.Password,
		Captcha:  param.Captcha,
		Name:     param.Name,
		Email:    param.Email,
	}
	_, err := userServiceClient.Register(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	//返回成功的响应码
	res.ResponseSuccess(c, nil)
}

func (u *HandlerUser) userLogin(c *gin.Context) {
	//1.将参数绑定到结构体中
	//2.调用grpc服务的用户登录服务
	//3.grpc服务中，先查询用户是否存在
	//4.获取用户信息.
	//5.返回响应

	res := common.NewResponseData()
	var param model.ParamLogin

	if err := c.ShouldBind(&param); err != nil {
		zap.L().Error("register with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			res.ResponseError(c, common.CodeInvalidParams)
			return
		}
		res.ResponseErrorWithMsg(c, common.CodeInvalidParams, errs.Translate(common.Trans))
		return
	}

}
