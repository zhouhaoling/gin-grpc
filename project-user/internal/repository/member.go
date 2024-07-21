package repository

import (
	"context"
	"time"

	"test.com/project-user/internal/repository/dao"

	"go.uber.org/zap"
	"test.com/common/errs"
	"test.com/project-grpc/user_grpc"
	"test.com/project-user/config"
	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repo"
	"test.com/project-user/internal/repository/database"
)

type MemberRepository struct {
	member repo.MemberRepo
	org    repo.OrganizationRepo
	adr    repo.AddressRepo
	dt     repo.DingTalkRepo
}

func NewMemberRepository() *MemberRepository {
	return &MemberRepository{
		member: dao.NewMemberDao(),
		org:    dao.NewOrganizationDao(),
		adr:    dao.NewAddressDao(),
		dt:     dao.NewDingTalkDao(),
	}
}

func (repo *MemberRepository) FindMemberByMobile(ctx context.Context, mobile string) (bool, error) {
	return repo.member.GetMemberByMobile(ctx, mobile)
}

// CreateMember 创建成员
func (repo *MemberRepository) CreateMember(conn database.DBConn, ctx context.Context, mid int64, pwd []byte, msg *user_grpc.RegisterRequest) (*model.Member, error) {
	member := &model.Member{
		Account:       msg.Name,
		Email:         msg.Email,
		Password:      string(pwd),
		Mobile:        msg.Mobile,
		Name:          msg.Name,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: time.Now().UnixMilli(),
		Status:        config.Normal,
		Avatar:        config.DefaultAvatar,
	}
	member.MId = mid
	if err := repo.member.InsertMember(conn, ctx, member); err != nil {
		return nil, err
	}
	address := &model.Address{}
	address.MId = mid
	if err := repo.adr.InsertAddress(conn, ctx, address); err != nil {
		return nil, err
	}
	return member, nil
}

// IsRegisterMemberExist 用于判断该注册用户的手机号码和邮箱是否已经存在
func (repo *MemberRepository) IsRegisterMemberExist(ctx context.Context, msg *user_grpc.RegisterRequest) error {
	exist, err := repo.member.GetMemberByEmail(ctx, msg.Email)
	if err != nil {
		zap.L().Error("register get email error", zap.Error(err))
		return errs.GrpcError(model.MySQLError)
	}
	if exist {
		zap.L().Warn("email is exist, register failed")
		return errs.GrpcError(model.EmailExist)
	}
	exist, err = repo.member.GetMemberByMobile(ctx, msg.Mobile)
	if err != nil {
		zap.L().Error("register get mobile error", zap.Error(err))
		return errs.GrpcError(model.MySQLError)
	}
	if exist {
		zap.L().Warn("mobile is exist, register failed")
		return errs.GrpcError(model.MobileExist)
	}
	return nil
}

func (repo *MemberRepository) FindMemberByAccount(ctx context.Context, account string) (model.Member, error) {
	return repo.member.SelectMemberByAccount(ctx, account)
}

func (repo *MemberRepository) FindMemberByMId(ctx context.Context, mid int64) (*model.Member, error) {
	return repo.member.SelectMemberByMemId(ctx, mid)
}

func (repo *MemberRepository) FindMemberByMidAndNameAndOrgId(ctx context.Context, mid int64, name string, orgId int64) (bool, error) {
	exist, err := repo.member.SelectMemberByMemIdAndName(ctx, mid, name)
	if err != nil {
		return false, err
	}
	if exist == false {
		return exist, nil
	}
	//从数据库中查询orgId
	exist, err = repo.org.SelectOrganizationByOrgId(ctx, orgId)
	if err != nil {
		return false, err
	}
	if exist == false {
		return false, nil
	}
	return true, nil
}

func (repo *MemberRepository) FindMemberByMIds(ctx context.Context, ids []int64) ([]*model.Member, error) {
	return repo.member.SelectMemberByMemIds(ctx, ids)
}
