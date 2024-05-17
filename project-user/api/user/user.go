package user

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"test.com/project-user/internal/dao"
	"test.com/project-user/internal/repo"

	"test.com/project-user/tools"

	"github.com/gin-gonic/gin"
	"test.com/common"
)

type HandlerUser struct {
	cache repo.Cache
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{
		cache: dao.RC,
	}
}

func (u *HandlerUser) getCaptcha(c *gin.Context) {
	rsp := common.NewResponseData()

	//1.获取参数
	mobile := c.PostForm("mobile")
	//2.校验参数
	if !common.VerifyModel(mobile) {
		rsp.ResponseError(c, common.CodeNoLegalMobile)
		return
	}
	//3.生成验证码
	code := tools.GetVerifyCode()
	fmt.Println("code:", code)
	//4.调用短信验证平台
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("调用短信验证平台成功,发送短信")
		zap.L().Error("调用失败")
		zap.L().Warn("警告")
		//redis 假设后续缓存可能存在mysql中，也可能存在mongo中，或者memcache中
		//存储验证码到redis中,并设置过期时间
		key := "register_" + mobile
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := u.cache.Put(ctx, key, code, 15*time.Minute)
		if err != nil {
			zap.L().Error("验证码存入redis出错, err:", zap.Error(err))
			return
		}
		zap.L().Info("验证码存入redis成功")
	}()
	rsp.ResponseSuccess(c, code)
}
