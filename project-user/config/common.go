package config

import "time"

// Redis key
var (
	RegisterMobileCacheKey = "register_"
	LoginMobileCacheKey    = "login_"
)

var (
	MaxLifetime = time.Hour
)
