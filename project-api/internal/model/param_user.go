package model

// ParamLogin 登录请求参数
type ParamLogin struct {
}

// ParamRegister 注册请求参数
type ParamRegister struct {
	Mobile     string `json:"mobile"`      //手机
	Email      string `json:"email"`       //邮箱
	Name       string `json:"name"`        //姓名
	Password   string `json:"password"`    //密码
	RePassword string `json:"re_password"` //确认密码
	Captcha    string `json:"captcha"`     //验证码
}
