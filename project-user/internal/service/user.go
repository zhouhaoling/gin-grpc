package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"test.com/project-user/internal/domain"

	pg "test.com/project-grpc/project_grpc"
	"test.com/project-user/internal/rpc"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"test.com/common"
	"test.com/common/encrypts"
	"test.com/common/errs"
	"test.com/common/jwts"
	"test.com/common/tms"
	"test.com/project-grpc/user_grpc"
	"test.com/project-user/config"
	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repo"
	"test.com/project-user/internal/repository"
	"test.com/project-user/internal/repository/database"
	"test.com/project-user/pkg/snowflake"
	"test.com/project-user/tools"
)

type UserService struct {
	user_grpc.UnimplementedLoginServiceServer
	user_grpc.UnimplementedUserServiceServer
	cache               repo.Cache
	mrepo               *repository.MemberRepository
	orepo               *repository.OrganizationRepository
	arepo               *repository.AddressRepository
	drepo               *repository.DingTalkRepository
	tran                database.Transaction
	memberAccountDomain *domain.MemberAccountDomain
}

func NewUserService(cache repo.Cache) *UserService {
	return &UserService{
		cache:               cache,
		mrepo:               repository.NewMemberRepository(),
		orepo:               repository.NewOrganizationRepository(),
		tran:                repository.NewTransaction(),
		arepo:               repository.NewAddressRepository(),
		drepo:               repository.NewDingTalkRepository(),
		memberAccountDomain: domain.NewMemberAccountDomain(),
	}
}

// Register 用于注册
func (svc *UserService) Register(c context.Context, msg *user_grpc.RegisterRequest) (*user_grpc.RegisterResponse, error) {
	//grpc注册服务中的步骤
	//1.获取redis中存储的验证码并对比
	//fmt.Println("用户密码：", msg.Password)
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
		organization, err := svc.orepo.CreateOrganization(conn, ctx, member)
		if err != nil {
			zap.L().Error("register organization failed", zap.Error(err))
			return errs.GrpcError(model.MySQLError)
		}
		oCode := encrypts.EncryptInt64NoErr(organization.Id)
		//插入项目模板信息，项目成员级别信息
		msg := &pg.ProjectRequest{
			MemberId:         mid,
			OrganizationCode: oCode,
			Name:             member.Name,
			Mobile:           member.Mobile,
			Email:            member.Email,
			Description:      member.Description,
			Avatar:           member.Avatar,
		}
		_, err = rpc.ProjectServiceClient.CopyProjectTemplate(ctx, msg)
		if err != nil {
			return err
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
	fmt.Println("msg.Password", msg.Password)
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
		Id:       member.Id,
		Mid:      member.MId,
		Name:     member.Name,
		Mobile:   member.Mobile,
		Realname: member.RealName,
		Account:  member.Account,
		Email:    member.Email,
		Status:   int32(member.Status),
		Address:  addr.Address,
		Province: int32(addr.Province),
		City:     int32(addr.City),
		Area:     int32(addr.Area),
		Avatar:   member.Avatar,
	}
	memberResp.Code, _ = encrypts.EncryptInt64(member.MId, config.AESKey)
	memberResp.LastLoginTime = tms.FormatByMill(member.LastLoginTime)
	memberResp.CreateTime = tms.FormatByMill(member.CreateTime)
	fmt.Println("memberResp:", memberResp)
	//查询组织表信息
	orgs, err := svc.orepo.FindOrganizationListByMId(ctx, member.MId)
	if err != nil {
		zap.L().Error("login fail, select organization failed, error:", zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.GrpcError(model.OrgNotExist)
		}
		return nil, errs.GrpcError(model.MySQLError)
	}
	//fmt.Println("执行到此处")
	var orgResp []*user_grpc.OrganizationResponse
	err = copier.Copy(&orgResp, orgs)
	if err != nil {
		zap.L().Error("copy failed, error:", zap.Error(err))
	}
	for _, v := range orgResp {
		//把项目id信息加密，返回给前端
		v.Code, _ = encrypts.EncryptInt64(v.Id, config.AESKey)
		v.OwnerCode = memberResp.Code
		v.CreateTime = tms.FormatByMill(model.ToMap(orgs)[v.Id].CreateTime)
	}
	fmt.Println("orgResp", orgResp)
	orgsCode, _ := encrypts.EncryptInt64(orgs[0].Id, config.AESKey)
	//jwt生成
	jwtToken, err := jwts.CreateToken(member.MId, member.Name, orgsCode, msg.Ip)
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
	//修改账号最后一次登录时间
	err = svc.memberAccountDomain.UpdateMemberAccountByLastTime(ctx, member.MId)
	if err != nil {
		return nil, errs.GrpcError(model.MySQLError)
	}
	//存入缓存中
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		marsh, _ := json.Marshal(member)
		midStr := strconv.FormatInt(member.MId, 10)
		svc.cache.Put(ctx, config.MemberCacheKey+"::"+midStr, string(marsh), config.CacheLifeTime)
		orgJson, _ := json.Marshal(orgs)
		svc.cache.Put(ctx, config.OrganCacheKey+"::"+midStr, string(orgJson), config.CacheLifeTime)
	}()
	return &user_grpc.LoginResponse{
		Member:           memberResp,
		OrganizationList: orgResp,
		TokenList:        tl,
	}, err
}

func (svc *UserService) MyOrgList(ctx context.Context, msg *user_grpc.UserRequest) (*user_grpc.OrgListResponse, error) {
	mid := msg.MemberId

	orgs, err := svc.orepo.FindOrganizationListByMId(ctx, mid)
	if err != nil {
		zap.L().Error("MyOrgList method select orgs by mid failed, error:", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	var orgResp []*user_grpc.OrganizationResponse
	err = copier.Copy(&orgResp, orgs)
	if err != nil {
		zap.L().Error("copy orgs => orgResp failed", zap.Error(err))
		return nil, errs.GrpcError(model.ErrorServerBusy)
	}
	for _, v := range orgResp {
		v.Code, _ = encrypts.EncryptInt64(v.Id, config.AESKey)
	}
	return &user_grpc.OrgListResponse{OrganizationList: orgResp}, nil
}

func (svc *UserService) FindMemberByMemId(ctx context.Context, msg *user_grpc.UserRequest) (*user_grpc.MemberResponse, error) {
	mid := msg.MemberId
	member, err := svc.mrepo.FindMemberByMId(ctx, mid)
	if err != nil {
		zap.L().Error("FindMemberByMemId FindMemberByMId failed", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}

	return &user_grpc.MemberResponse{
		Id:     member.Id,
		Mid:    member.MId,
		Name:   member.Name,
		Avatar: member.Avatar,
	}, nil
}

func (svc *UserService) TokenVerify(ctx context.Context, msg *user_grpc.TokenRequest) (*user_grpc.TokenMessage, error) {
	//msg中有member_id, member_name, organCode
	//1.从缓存中查询数据，查不到再到数据库中查
	//2.对比数据
	//3.返回响应
	midStr := strconv.FormatInt(msg.MemberId, 10)
	orgsCode, _ := encrypts.Decrypt(msg.OrganCode, config.AESKey) //org[0].Id
	orgId, _ := strconv.ParseInt(orgsCode, 10, 64)
	c, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	//TODO 从缓存中查，查询不到直接返回
	if err := svc.tokenVerifyRedis(c, midStr, orgId, msg); err != nil {
		//exist, err := svc.mrepo.FindMemberByMidAndNameAndOrgId(c, msg.MemberId, msg.MemberName, orgId)
		//if err != nil || exist == false {
		//	zap.L().Error("数据库中没有该记录")
		//	return nil, errs.GrpcError(model.NoLogin)
		//}
		return nil, err
	}
	return &user_grpc.TokenMessage{
		Flag: true,
	}, nil
}

// tokenVerifyRedis 从缓存中查询数据
func (svc *UserService) tokenVerifyRedis(ctx context.Context, mid string, pid int64, msg *user_grpc.TokenRequest) error {
	memberJson, err := svc.cache.Get(ctx, config.MemberCacheKey+"::"+mid)
	if err != nil {
		zap.L().Error("user tokenVerifyRedis Get() memberJson failed", zap.Error(err))
		return errs.GrpcError(model.RedisError)
	}
	if memberJson == "" {
		zap.L().Error("memberJson expire")
		return errs.GrpcError(model.NoLogin)
	}
	member := &model.Member{}
	json.Unmarshal([]byte(memberJson), member)
	if msg.MemberName != member.Name || msg.MemberId != member.MId {
		zap.L().Error("member name not match")
		return errs.GrpcError(model.NoLogin)
	}
	orgJson, err := svc.cache.Get(ctx, config.OrganCacheKey+"::"+mid)
	if err != nil {
		zap.L().Error("user TokenVerifyRedis Get orgJson failed", zap.Error(err))
		return errs.GrpcError(model.RedisError)
	}
	if orgJson == "" {
		zap.L().Error("orgJson expire")
		return errs.GrpcError(model.NoLogin)
	}
	var org []*model.Organization
	json.Unmarshal([]byte(orgJson), &org)
	if pid != org[0].Id {
		zap.L().Error("organization id not match")
		return errs.GrpcError(model.NoLogin)
	}
	return nil
}

func (svc *UserService) FindMemberInfoByMIds(ctxs context.Context, msg *user_grpc.UserRequest) (*user_grpc.MemberListResponse, error) {
	//fmt.Println("调用FindMemberInfoByMIds成功")
	ctx, cancel := context.WithTimeout(ctxs, 4*time.Second)
	defer cancel()
	memberList, err := svc.mrepo.FindMemberByMIds(ctx, msg.MIds)
	if err != nil {
		zap.L().Error("user FindMemberInfoByMIds FindMemberByMIds() failed", zap.Error(err))
		return nil, errs.GrpcError(model.MySQLError)
	}
	if memberList == nil || len(memberList) <= 0 {
		return &user_grpc.MemberListResponse{List: nil}, nil
	}
	mMap := make(map[int64]*model.Member)
	for _, v := range memberList {
		mMap[v.Id] = v
	}
	var memMsgs []*user_grpc.MemberResponse
	copier.Copy(&memMsgs, memberList)
	for _, v := range memMsgs {
		m := mMap[v.Id]
		//fmt.Println("v.id:", v.Id)
		v.CreateTime = tms.FormatByMill(m.CreateTime)
		v.Code = encrypts.EncryptInt64NoErr(m.MId)
		v.Mid = m.MId
	}
	return &user_grpc.MemberListResponse{
		List: memMsgs,
	}, nil
}
