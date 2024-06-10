package repository

import (
	"context"
	"time"

	"test.com/project-user/internal/repository/database"

	"test.com/project-user/config"

	"go.uber.org/zap"
	"test.com/common/errs"
	"test.com/project-grpc/user_grpc"
	"test.com/project-user/internal/model"
	"test.com/project-user/internal/repository/dao/mysql"
)

type MemberRepository struct {
	member *mysql.MemberDao
	org    *mysql.OrganizationDao
	adr    *mysql.AddressDao
	dt     *mysql.DingTalkDao
}

func NewMemberRepository() *MemberRepository {
	return &MemberRepository{
		member: mysql.NewMemberDao(),
		org:    mysql.NewOrganizationDao(),
		adr:    mysql.NewAddressDao(),
		dt:     mysql.NewDingTalkDao(),
	}
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
	dingtalk := &model.DingTalk{}
	dingtalk.MId = mid
	if err := repo.dt.InsertDingTalk(conn, ctx, dingtalk); err != nil {
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
