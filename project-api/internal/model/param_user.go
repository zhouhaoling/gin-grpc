package model

import (
	"errors"

	"test.com/common"
)

// ParamLogin 登录请求参数
type ParamLogin struct {
}

// ParamRegister 注册请求参数
type ParamRegister struct {
	Mobile     string `json:"mobile" form:"mobile" binding:"required"`                            //手机
	Email      string `json:"email" form:"email" binding:"required"`                              //邮箱
	Name       string `json:"name" form:"name" binding:"required"`                                //姓名
	Password   string `json:"password" form:"password" binding:"required"`                        //密码
	RePassword string `json:"re_password" form:"re_password" binding:"required,eqfield=Password"` //确认密码
	Captcha    string `json:"captcha" form:"captcha" binding:"required"`                          //验证码
}

func (p ParamRegister) Verify() error {
	if !common.VerifyModel(p.Mobile) {
		return errors.New("手机格式不正确")
	}
	if !common.VerifyCode(p.Captcha) {
		return errors.New("验证码格式不正确")
	}
	if !common.VerifyEmail(p.Email) {
		return errors.New("邮箱格式不正确")
	}
	return nil
}
