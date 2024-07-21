package config

import "time"

const (
	CtxMemberIDKey       = "member_id"
	CtxMemberNameKey     = "member_name"
	CtxOrganizationIDKey = "organization_id"

	AESKey     = "abcdefgehjhijkmlkjjwwoew"
	CtxTimeOut = 3 * time.Second
)

var MemberPaths = []string{
	"/project/login/getCaptcha",
	"/project/login",
	"/project/login/register",
}
