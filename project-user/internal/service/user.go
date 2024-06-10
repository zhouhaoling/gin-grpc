package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"test.com/project-user/internal/repository/database"

	"github.com/go-redis/redis/v8"

	"test.com/project-user/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"test.com/project-user/pkg/snowflake"

	"test.com/project-user/config"

	"test.com/project-grpc/user_grpc"

	"test.com/common/errs"

	"test.com/project-user/internal/model"

	"test.com/project-user/internal/repo"

	"go.uber.org/zap"
	"test.com/common"
	"test.com/project-user/tools"
)

type UserService struct {
	user_grpc.UnimplementedLoginServiceServer
	user_grpc.UnimplementedUserServiceServer
	cache repo.Cache
	mrepo *repository.MemberRepository
	orepo *repository.OrganizationRepository
	arepo *repository.AddressRepository
	drepo *repository.DingTalkRepository
	tran  database.Transaction
}

func NewUserService(cache repo.Cache) *UserService {
	return &UserService{
		cache: cache,
		mrepo: repository.NewMemberRepository(),
		orepo: repository.NewOrganizationRepository(),
		tran:  repository.NewTransaction(),
		arepo: repository.NewAddressRepository(),
		drepo: repository.NewDingTalkRepository(),
	}
}

// Register 用于注册
func (svc *UserService) Register(c context.Context, msg *user_grpc.RegisterRequest) (*user_grpc.RegisterResponse, error) {
	//grpc注册服务中的步骤
	//1.获取redis中存储的验证码并对比
	ctx := context.Background()
	key := config.RegisterMobileCacheKey + msg.Mobile
	captcha, err := svc.cache.Get(ctx, key)
	if errors.Is(err, redis.Nil) {
		zap.L().Warn("验证码已过期")
		return nil, errs.GrpcError(model.CaptchaNotExist)
	}
	if err != nil {
		zap.L().Error("register redis get error", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if captcha != msg.Captcha {
		zap.L().Warn("captcha not equal to msg.Captcha")
		return nil, errs.GrpcError(model.ErrorCaptcha)
	}
	//2.校验业务逻辑(邮箱是否被注册，手机是否被注册）
	err = svc.mrepo.IsRegisterMemberExist(ctx, msg)
	if err != nil {
		return nil, err
	}
	//3.执行业务,生成uuid，并将用户信息存入mysql中的organization表和member表
	mid := snowflake.GenID()
	hash, err := bcrypt.GenerateFromPassword([]byte(msg.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Warn("使用bcrypt加密失败", zap.Error(err))
		return nil, errs.GrpcError(model.ErrorServerBusy)
	}
	//开启事物
	err = svc.tran.Action(func(conn database.DBConn) error {
		//插入用户信息
		member, err := svc.mrepo.CreateMember(conn, ctx, mid, hash, msg)
		if err != nil {
			zap.L().Error("register member failed", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		//插入
		//插入组织信息
		err = svc.orepo.CreateOrganization(conn, ctx, member)
		if err != nil {
			zap.L().Error("register organization failed", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		return nil
	})

	//4.返回响应
	return &user_grpc.RegisterResponse{}, err
}

func (svc *UserService) GetCaptcha(c context.Context, msg *user_grpc.CaptchaRequest) (*user_grpc.CaptchaResponse, error) {
	//1.获取参数
	mobile := msg.Mobile
	//2.校验参数
	if !common.VerifyModel(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	//3.生成验证码
	code := tools.GetVerifyCode()
	fmt.Println("code:", code)
	//4.调用短信验证平台
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("调用短信验证平台成功,发送短信")
		//redis存储 假设后续缓存可能存在mysql中，也可能存在mongo中，或者memcache中
		//存储验证码到redis中,并设置过期时间
		key := config.RegisterMobileCacheKey + mobile
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := svc.cache.Put(ctx, key, code, 15*time.Minute)
		if err != nil {
			zap.L().Error("验证码存入redis出错, err:", zap.Error(err))
			return
		}

		zap.L().Info("验证码存入redis成功")
	}()
	return &user_grpc.CaptchaResponse{
		Code: code,
	}, nil
}
