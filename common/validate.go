package common

import (
	"github.com/dlclark/regexp2"
)

func VerifyModel(mobile string) bool {
	pattern := `^1[3-9]\d{9}$`
	reg := regexp2.MustCompile(pattern, regexp2.None)
	match, err := reg.MatchString(mobile)
	if err != nil {
		return false
	}
	return match
}
