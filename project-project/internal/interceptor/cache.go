package interceptor

import (
	"context"
	"encoding/json"
	"time"

	pg "test.com/project-grpc/project_grpc"

	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/repository/dao"

	"go.uber.org/zap"
	"test.com/common/encrypts"

	"google.golang.org/grpc"
)

// CacheInterceptor 接口缓存拦截器 // TODO 未实现日志拦截器，打印参数内容值，请求时间等
type CacheInterceptor struct {
	cache    repo.Cache
	cacheMap map[string]interface{}
}

func NewCacheInterceptor() *CacheInterceptor {
	cacheMap := make(map[string]interface{})
	cacheMap["/proto.ProjectService/FindProjectByMemberId"] = &pg.ProjectResponse{}
	return &CacheInterceptor{
		cache:    dao.NewRedisCache(),
		cacheMap: cacheMap,
	}
}

func (c *CacheInterceptor) Cache() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		respType := c.cacheMap[info.FullMethod]
		if respType == nil {
			return handler(ctx, req)
		}
		//先查询是否有缓存，有的话， 直接返回，无，再请求，存入缓存
		con, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		marshal, _ := json.Marshal(req)
		cacheKey := encrypts.Md5(string(marshal))
		respJson, _ := c.cache.Get(con, info.FullMethod+"::"+cacheKey)
		if respJson != "" {
			json.Unmarshal([]byte(respJson), &respType)
			zap.L().Info(info.FullMethod + "查询缓存成功")
			return respType, nil
		}
		resp, err = handler(ctx, req)
		bytes, _ := json.Marshal(resp)
		c.cache.Put(con, info.FullMethod+"::"+cacheKey, string(bytes), 5*time.Minute)
		zap.L().Info(info.FullMethod + "存入缓存成功")
		return
	})
}

func (c *CacheInterceptor) CacheInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		respType := c.cacheMap[info.FullMethod]
		if respType == nil {
			return handler(ctx, req)
		}
		//先查询是否有缓存，有的话， 直接返回，无，再请求，存入缓存
		con, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		marshal, _ := json.Marshal(req)
		cacheKey := encrypts.Md5(string(marshal))
		respJson, _ := c.cache.Get(con, info.FullMethod+"::"+cacheKey)
		if respJson != "" {
			json.Unmarshal([]byte(respJson), &respType)
			zap.L().Info(info.FullMethod + "查询缓存成功")
			return respType, nil
		}
		resp, err = handler(ctx, req)
		bytes, _ := json.Marshal(resp)
		c.cache.Put(con, info.FullMethod+"::"+cacheKey, string(bytes), 5*time.Minute)
		zap.L().Info(info.FullMethod + "存入缓存成功")
		return
	}
}
