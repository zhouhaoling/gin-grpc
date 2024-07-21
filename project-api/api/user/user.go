package user

import (
	"context"
	"fmt"
	"time"

	"test.com/project-api/internal/tool"

	"test.com/project-api/api/rpc"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"test.com/common"
	"test.com/common/errs"
	"test.com/project-api/config"
	"test.com/project-api/internal/model"
	ug "test.com/project-grpc/user_grpc"
)

type HandlerUser struct {
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}

// getCaptcha 获取验证码
func (u *HandlerUser) getCaptcha(c *gin.Context) {
	res := common.NewResponseData()
	mobile := c.PostForm("mobile")
	codeType := c.PostForm("type")
	fmt.Println("type:", codeType)
	if flag := common.VerifyModel(mobile); flag != true {
		res.ResponseErrorWithMsg(c, common.CodeInvalidParams, "手机号格式不正确")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var response *ug.CaptchaResponse
	var err error
	switch codeType {
	case "0":
		response, err = rpc.LoginServiceClient.GetRegisterCaptcha(ctx, &ug.CaptchaRequest{Mobile: mobile})
	case "1":
		response, err = rpc.LoginServiceClient.GetLoginCaptcha(ctx, &ug.CaptchaRequest{Mobile: mobile})
	default:
		res.ResponseError(c, common.CodeServerBusy)
		return
	}

	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	//返回成功的响应码
	res.ResponseSuccess(c, response.Code)
}

// userRegister 用户注册
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
	_, err := rpc.UsersServiceClient.Register(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	//返回成功的响应码
	res.ResponseSuccess(c, nil)
}

// userLogin 用户登录
func (u *HandlerUser) userLogin(c *gin.Context) {
	//1.将参数绑定到结构体中
	//2.调用grpc服务的用户登录服务
	//3.返回响应

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
	//调用grpc服务，完成登录
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ip := tool.GetIP(c)
	msg := &ug.LoginRequest{
		Account:  param.Account,
		Password: param.Password,
		Ip:       ip,
	}
	response, err := rpc.UsersServiceClient.Login(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	loginResp := &model.LoginResp{}
	err = copier.Copy(loginResp, response)
	if err != nil {
		res.ResponseError(c, common.CodeServerBusy)
		return
	}
	//返回响应
	res.ResponseSuccess(c, loginResp)
}

// myOrgList 我的组织列表
func (u *HandlerUser) myOrgList(c *gin.Context) {
	//1.通过上下文获取用户id
	//2.调用grpc服务
	//3.返回响应
	res := common.NewResponseData()
	midAny, exist := c.Get(config.CtxMemberIDKey)
	if !exist {
		zap.L().Error("member not login")
		res.ResponseError(c, common.CodeServerBusy)
		return
	}
	mid := midAny.(int64)
	fmt.Println("member_id", mid)
	myOrgList, err := rpc.UsersServiceClient.MyOrgList(context.Background(), &ug.UserRequest{MemberId: mid})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		res.ResponseErrorWithMsg(c, code, msg)
		return
	}
	if myOrgList == nil {
		res.ResponseSuccess(c, []*model.OrganizationList{})
		return
	}
	ol := myOrgList.OrganizationList
	var orgs []*model.OrganizationList
	copier.Copy(&orgs, ol)
	res.ResponseSuccess(c, orgs)
}
