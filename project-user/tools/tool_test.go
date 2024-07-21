package tools

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestGetVerifyCode(t *testing.T) {
	code := GetVerifyCode1()
	fmt.Println("code:", code)
	pwd, _ := bcrypt.GenerateFromPassword([]byte("123456zhl"), bcrypt.DefaultCost)
	fmt.Printf("pwd:%s。密码", pwd)
}

// 该方法有5位数验证码生成。
func GetVerifyCode1() string {
	//为随机数生成器设置种子，使用当前时间
	rand.New(rand.NewSource(time.Now().UnixNano()))

	//生成一个0到999999之间的随机数，并取最后6位作为验证码
	get := rand.Intn(1000000)
	code := get % 1000000 //取最后6位
	fmt.Println("随机生成的6位数验证码:", code)
	//转为字符串
	codeStr := strconv.Itoa(code)
	//如果codeStr的长度为5，则重新生成
	if len(codeStr) != 6 {
		return GetVerifyCode()
	}
	return codeStr
}
