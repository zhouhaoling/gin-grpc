package domain

import (
	"context"

	"go.uber.org/zap"

	"test.com/project-user/internal/repo"
	"test.com/project-user/internal/repository/dao"
)

type MemberAccountDomain struct {
	accountRepo repo.MemberAccountRepo
}

func (d *MemberAccountDomain) UpdateMemberAccountByLastTime(ctx context.Context, mid int64) error {
	err := d.accountRepo.UpdateMemberAccountByLastTime(ctx, mid)
	if err != nil {
		zap.L().Error("domain MemberAccountDomain UpdateMemberAccountByLastTime().accountRepo.UpdateMemberAccountByLastTime() failed", zap.Error(err))
		return err
	}
	return nil
}

func NewMemberAccountDomain() *MemberAccountDomain {
	return &MemberAccountDomain{
		accountRepo: dao.NewMemberAccountDao(),
	}
}
