package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	ug "test.com/project-user/user_grpc"

	"test.com/project-user/internal/repo"

	"go.uber.org/zap"
	"test.com/common"
	"test.com/project-user/tools"
)

type UserService struct {
	ug.UnimplementedLoginServiceServer
	cache repo.Cache
}

func NewUserService(cache repo.Cache) *UserService {
	return &UserService{
		cache: cache,
	}
}

func (svc *UserService) GetCaptcha(ctx context.Context, msg *ug.CaptchaRequest) (*ug.CaptchaResponse, error) {
	//1.获取参数
	mobile := msg.Mobile
	//2.校验参数
	if !common.VerifyModel(mobile) {
		return nil, errors.New("手机号不合法")
	}
	//3.生成验证码
	code := tools.GetVerifyCode()
	fmt.Println("code:", code)
	//4.调用短信验证平台
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("调用短信验证平台成功,发送短信")
		//redis 假设后续缓存可能存在mysql中，也可能存在mongo中，或者memcache中
		//存储验证码到redis中,并设置过期时间
		key := "register_" + mobile
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		err := svc.cache.Put(ctx, key, code, 15*time.Minute)
		if err != nil {
			zap.L().Error("验证码存入redis出错, err:", zap.Error(err))
			return
		}
		zap.L().Info("验证码存入redis成功")
	}()
	return &ug.CaptchaResponse{
		Code: code,
	}, nil
}
