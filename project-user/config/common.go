package config

import "time"

// Redis key
var (
	RegisterMobileCacheKey = "register_"
	LoginMobileCacheKey    = "login_"
)

var (
	MaxLifetime   = time.Hour
	DefaultAvatar = "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5"
)

var (
	Normal         = 1
	Personal int32 = 1
)

var (
	PersonalOrganization = "个人组织"
)
