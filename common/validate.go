package common

import "github.com/dlclark/regexp2"

// VerifyModel 校验11位手机号
func VerifyModel(mobile string) bool {
	pattern := `^1[3-9]\d{9}$`
	reg := regexp2.MustCompile(pattern, regexp2.None)
	match, err := reg.MatchString(mobile)
	if err != nil {
		return false
	}
	return match
}

// VerifyCode 校验6位验证码
func VerifyCode(code string) bool {
	pattern := `\d{6}$`
	reg := regexp2.MustCompile(pattern, regexp2.None)
	match, err := reg.MatchString(code)
	if err != nil {
		return false
	}
	return match
}

// VerifyEmail 校验邮箱
func VerifyEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	reg := regexp2.MustCompile(pattern, regexp2.Debug)
	match, err := reg.MatchString(email)
	if err != nil {
		return false
	}
	return match
}

// VerifyPassword 校验密码
func VerifyPassword(password string) bool {
	pattern := `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{6,}$`
	reg := regexp2.MustCompile(pattern, regexp2.None)
	match, err := reg.MatchString(password)
	if err != nil {
		return false
	}
	return match
}
