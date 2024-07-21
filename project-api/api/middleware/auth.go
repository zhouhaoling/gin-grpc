package middleware

import (
	"context"
	"net/http"
	"strings"

	"test.com/project-api/internal/tool"

	rpc "test.com/project-api/api/rpc"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"test.com/common/jwts"
	"test.com/project-api/config"
	ug "test.com/project-grpc/user_grpc"
)

type JwtMiddlewareBuilder struct {
	paths []string
}

func NewJwtMiddlewareBuilder() *JwtMiddlewareBuilder {
	return &JwtMiddlewareBuilder{}
}

func (jmb *JwtMiddlewareBuilder) IgnorePath(path string) *JwtMiddlewareBuilder {
	jmb.paths = append(jmb.paths, path)
	return jmb
}
func (jmb *JwtMiddlewareBuilder) IgnorePaths(paths []string) *JwtMiddlewareBuilder {
	for _, value := range paths {
		jmb.paths = append(jmb.paths, value)
	}
	return jmb
}

func (jmb *JwtMiddlewareBuilder) Build() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		//1.判断请求是否在白名单
		//2.从header中获取token
		//3.处理结果，将member_id设置到gin的上下文中,失败则返回未登录
		for _, path := range jmb.paths {
			if path == c.Request.URL.Path {
				return
			}
		}
		tokenHeader := c.GetHeader("Authorization")
		if tokenHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ip := tool.GetIP(c)
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 && segs[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims, err := jwts.ParseToken(tokenStr)
		if err != nil {
			zap.L().Error("token解析失败", zap.Error(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims == nil || claims.MemberId == 0 || claims.Ip != ip {
			zap.L().Error("claims数据为空", zap.Error(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//TODO 未补充：先去查询node表，确认不使用登录控制的接口，就不做登录认证
		_, err = rpc.UsersServiceClient.TokenVerify(context.Background(), &ug.TokenRequest{
			MemberId:   claims.MemberId,
			MemberName: claims.MemberName,
			OrganCode:  claims.OrganizationId,
		})
		if err != nil {
			zap.L().Error("token验证失败", zap.Error(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//浏览器不同，user-agent不同，此时返回
		//if claims.UserAgent != ctx.Request.UserAgent() {
		//	//监控这里
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		c.Set(config.CtxMemberIDKey, claims.MemberId)
		c.Set(config.CtxMemberNameKey, claims.MemberName)
		c.Set(config.CtxOrganizationIDKey, claims.OrganizationId)
		c.Next()
	}
}
