package domain

import (
	"context"
	"time"

	pu "test.com/project-grpc/user_grpc"
	"test.com/project-project/internal/rpc"
)

type UserRpcDomain struct {
	lc pu.UserServiceClient
}

func NewUserRpcDomain() *UserRpcDomain {
	return &UserRpcDomain{
		lc: rpc.UserServiceClient,
	}
}

func (d *UserRpcDomain) MemberList(mIdList []int64) ([]*pu.MemberResponse, map[int64]*pu.MemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	messageList, err := d.lc.FindMemberInfoByMIds(ctx, &pu.UserRequest{MIds: mIdList})
	mMap := make(map[int64]*pu.MemberResponse)
	for _, v := range messageList.List {
		mMap[v.Mid] = v
	}
	return messageList.List, mMap, err
}

func (d *UserRpcDomain) MemberInfo(mid int64) (*pu.MemberResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	message, err := d.lc.FindMemberByMemId(ctx, &pu.UserRequest{MemberId: mid})
	return message, err
}
