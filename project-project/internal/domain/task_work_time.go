package domain

import (
	"context"
	"time"

	"go.uber.org/zap"
	"test.com/common/errs"
	"test.com/project-project/internal/model"
	"test.com/project-project/internal/repository"
)

type TaskWorkTimeDomain struct {
	taskWorkTimeRepo *repository.TaskStagesRepository
	userRpcDomain    *UserRpcDomain
}

func NewTaskWorkTimeDomain() *TaskWorkTimeDomain {
	return &TaskWorkTimeDomain{
		taskWorkTimeRepo: repository.NewTaskStagesRepository(),
		userRpcDomain:    NewUserRpcDomain(),
	}
}

// TaskWorkTimeList 处理工时相关的逻辑
func (d *TaskWorkTimeDomain) TaskWorkTimeList(taskCode int64) ([]*model.TaskWorkTimeDisplay, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var list []*model.TaskWorkTime
	var err error
	list, err = d.taskWorkTimeRepo.FindWorkTimeList(c, taskCode)
	if err != nil {
		zap.L().Error("project task TaskWorkTimeList taskWorkTimeRepo.FindWorkTimeList error", zap.Error(err))
		return nil, model.MySQLError
	}
	if len(list) == 0 {
		return []*model.TaskWorkTimeDisplay{}, nil
	}
	var displayList []*model.TaskWorkTimeDisplay
	var mIdList []int64
	for _, v := range list {
		mIdList = append(mIdList, v.MemberCode)
	}
	_, mMap, err := d.userRpcDomain.MemberList(mIdList)
	if err != nil {
		return nil, errs.ToBError(err)
	}
	for _, v := range list {
		display := v.ToDisplay()
		message := mMap[v.MemberCode]
		m := model.Member{}
		m.Name = message.Name
		m.Id = message.Id
		m.Avatar = message.Avatar
		m.Code = message.Code
		display.Member = m
		displayList = append(displayList, display)
	}
	return displayList, nil
}
