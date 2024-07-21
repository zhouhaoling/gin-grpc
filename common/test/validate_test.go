package test

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"test.com/common"
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	//qq := "1735950045@qq.com"
	//flag := IsValidEmail(qq) true
	//com := "zhl19174683866@163.com"
	//flag := IsValidEmail(com) true
	//gmail := "zhl19174683866@gmail.com"
	//flag := IsValidEmail(gmail) // true
	//fmt.Println("flag:", flag)

	//test1 := "@qq.com" //false
	test2 := "125412@" //false
	flag := IsValidEmail(test2)
	fmt.Println("flag:", flag)
}

func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	reg := regexp2.MustCompile(pattern, regexp2.Debug)
	match, err := reg.MatchString(email)
	if err != nil {
		return false
	}
	return match
}

func TestIsValidMobile(t *testing.T) {
	mobile := "19174683866"
	flag := common.VerifyModel(mobile)
	fmt.Println("flag:", flag)
}
