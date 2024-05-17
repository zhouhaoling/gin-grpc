package tools

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GetVerifyCode() string {
	//为随机数生成器设置种子，使用当前时间
	rand.Seed(time.Now().UnixNano())

	//生成一个0到999999之间的随机数，并取最后6位作为验证码
	get := rand.Intn(1000000)
	code := get % 1000000 //取最后6位
	fmt.Println("随机生成的6位数验证码:", code)
	//转为字符串
	codeStr := strconv.Itoa(code)
	return codeStr
}
