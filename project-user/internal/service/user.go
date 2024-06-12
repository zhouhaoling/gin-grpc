package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"test.com/project-user/pkg/encrypts"

	"test.com/common/jwts"

	"github.com/jinzhu/copier"

	"gorm.io/gorm"

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
	pwd, err := bcrypt.GenerateFromPassword([]byte(msg.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Warn("使用bcrypt加密失败", zap.Error(err))
		return nil, errs.GrpcError(model.ErrorServerBusy)
	}
	//开启事物
	err = svc.tran.Action(func(conn database.DBConn) error {
		//插入用户信息
		member, err := svc.mrepo.CreateMember(conn, ctx, mid, pwd, msg)
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

func (svc *UserService) GetLoginCaptcha(c context.Context, msg *user_grpc.CaptchaRequest) (*user_grpc.CaptchaResponse, error) {
	//1.获取参数
	mobile := msg.Mobile
	//2.校验参数
	if !common.VerifyModel(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	//3.判断电话号码是否已被注册
	exist, err := svc.mrepo.FindMemberByMobile(context.Background(), mobile)
	if err != nil {
		zap.L().Error("find member mobile failed", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if !exist {
		zap.L().Warn("该手机号不存在")
		return nil, errs.GrpcError(model.MobileNotExist)
	}
	//4.生成验证码
	code := tools.GetVerifyCode()
	fmt.Println("code:", code)
	//4.调用短信验证平台
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("调用短信验证平台成功,发送短信")
		//redis存储 假设后续缓存可能存在mysql中，也可能存在mongo中，或者memcache中
		//存储验证码到redis中,并设置过期时间
		key := config.LoginMobileCacheKey + mobile
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

func (svc *UserService) GetRegisterCaptcha(c context.Context, msg *user_grpc.CaptchaRequest) (*user_grpc.CaptchaResponse, error) {
	//1.获取参数
	mobile := msg.Mobile
	//2.校验参数
	if !common.VerifyModel(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	//3.判断电话号码是否已被注册
	exist, err := svc.mrepo.FindMemberByMobile(context.Background(), mobile)
	if err != nil {
		zap.L().Error("find member mobile failed", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if exist {
		return nil, errs.GrpcError(model.MobileExist)
	}
	//4.生成验证码
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

func (svc *UserService) Login(c context.Context, msg *user_grpc.LoginRequest) (*user_grpc.LoginResponse, error) {
	//1.查询账号密码是否正确
	//2.获取用户信息，获取组织信息，生成token
	//3.返回响应

	ctx := context.Background()
	//查询账号是否存在
	member, err := svc.mrepo.FindMemberByAccount(ctx, msg.Account)
	if err != nil {
		zap.L().Error(" login fail, select member account failed, error:", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.GrpcError(model.AccountAndPwdError)
		}
		return nil, errs.GrpcError(model.MySQLError)
	}
	//验证密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(msg.Password))
	if err != nil {
		zap.L().Error("两个密码不一致,error:", zap.Error(err))
		return nil, errs.GrpcError(model.AccountAndPwdError)
	}
	//查询地址信息
	addr, err := svc.arepo.FindAddressByMId(ctx, member.MId)
	if err != nil {
		zap.L().Error("select address failed,error:", zap.Error(err))
	}
	memberResp := &user_grpc.MemberResponse{
		Id:            member.Id,
		Mid:           member.MId,
		Name:          member.Name,
		Mobile:        member.Mobile,
		Realname:      member.RealName,
		Account:       member.Account,
		Email:         member.Email,
		LastLoginTime: member.LastLoginTime,
		Status:        int32(member.Status),
		Address:       addr.Address,
		Province:      int32(addr.Province),
		City:          int32(addr.City),
		Area:          int32(addr.Area),
	}
	memberResp.Code, _ = encrypts.EncryptInt64(member.MId, config.AESKey)
	fmt.Println("memberResp:", memberResp)
	//查询组织表信息
	org, err := svc.orepo.FindOrganizationListByMId(ctx, member.MId)
	if err != nil {
		zap.L().Error("login fail, select organization failed, error:", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.GrpcError(model.OrgNotExist)
		}
		return nil, errs.GrpcError(model.MySQLError)
	}
	//fmt.Println("执行到此处")
	var orgResp []*user_grpc.OrganizationResponse
	err = copier.Copy(&orgResp, org)
	if err != nil {
		zap.L().Error("copy failed, error:", zap.Error(err))
	}
	for _, v := range orgResp {
		v.Code, _ = encrypts.EncryptInt64(v.Id, config.AESKey)
		v.Mbid, _ = encrypts.EncryptInt64(v.MemberId, config.AESKey)
	}
	//jwt生成
	jwtToken, err := jwts.CreateToken(member.MId)
	if err != nil {
		zap.L().Error("jwt生成失败,error:", zap.Error(err))
		return nil, errs.GrpcError(model.ErrorServerBusy)
	}
	tl := &user_grpc.TokenResponse{
		AccessToken:    jwtToken.AccessToken,
		RefreshToken:   jwtToken.RefreshToken,
		AccessTokenExp: jwtToken.AccessExp,
		TokenType:      config.AppConf.Jwt.TokenType,
	}
	return &user_grpc.LoginResponse{
		Member:           memberResp,
		OrganizationList: orgResp,
		TokenList:        tl,
	}, err
}
