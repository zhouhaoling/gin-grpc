package repo

import (
	"context"
	"time"
)

// Cache 用于实现缓存的接口类
type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
