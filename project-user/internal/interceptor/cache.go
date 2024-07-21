package interceptor

import (
	"context"
	"encoding/json"
	"time"

	"test.com/project-grpc/user_grpc"

	"go.uber.org/zap"
	"test.com/common/encrypts"

	"google.golang.org/grpc"
	"test.com/project-user/internal/repo"
	"test.com/project-user/internal/repository/dao/redis_dao"
)

// CacheInterceptor 接口缓存拦截器
type CacheInterceptor struct {
	cache    repo.Cache
	cacheMap map[string]interface{}
}

func NewCacheInterceptor() *CacheInterceptor {
	cacheMap := make(map[string]interface{})
	cacheMap["/proto.UserService/MyOrgList"] = &user_grpc.OrgListResponse{}
	cacheMap["/proto.UserService/FindMemberByMemId"] = &user_grpc.MemberResponse{}
	return &CacheInterceptor{
		cache:    redis_dao.RC,
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
