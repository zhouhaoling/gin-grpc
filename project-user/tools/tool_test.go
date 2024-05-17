package tools

import (
	"fmt"
	"testing"
)

func TestGetVerifyCode(t *testing.T) {
	code := GetVerifyCode()
	fmt.Println("code:", code)

}
