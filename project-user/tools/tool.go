package tools

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

//func GetVerifyCode() string {
//	//为随机数生成器设置种子，使用当前时间
//	rand.New(rand.NewSource(time.Now().UnixNano()))
//
//	//生成一个0到999999之间的随机数，并取最后6位作为验证码
//	get := rand.Intn(1000000)
//	code := get % 1000000 //取最后6位
//	fmt.Println("随机生成的6位数验证码:", code)
//	//转为字符串
//	codeStr := strconv.Itoa(code)
//	//如果codeStr的长度为5，则重新生成
//	if len(codeStr) == 5 {
//		return GetVerifyCode()
//	}
//	return codeStr
//}

func GetVerifyCode() string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(uint64(time.Now().UnixNano()))
	var length = 6
	var sb strings.Builder
	for i := 0; i < length; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
