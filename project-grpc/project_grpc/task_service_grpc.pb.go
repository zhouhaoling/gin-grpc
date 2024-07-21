// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.27.0
// source: task_service.proto

package project_grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	TaskService_TaskStages_FullMethodName        = "/proto.TaskService/TaskStages"
	TaskService_MemberProjectList_FullMethodName = "/proto.TaskService/MemberProjectList"
	TaskService_TaskList_FullMethodName          = "/proto.TaskService/TaskList"
	TaskService_SaveTask_FullMethodName          = "/proto.TaskService/SaveTask"
	TaskService_TaskSort_FullMethodName          = "/proto.TaskService/TaskSort"
	TaskService_MyTaskList_FullMethodName        = "/proto.TaskService/MyTaskList"
	TaskService_ReadTask_FullMethodName          = "/proto.TaskService/ReadTask"
	TaskService_ListTaskMember_FullMethodName    = "/proto.TaskService/ListTaskMember"
	TaskService_TaskLog_FullMethodName           = "/proto.TaskService/TaskLog"
	TaskService_TaskWorkTimeList_FullMethodName  = "/proto.TaskService/TaskWorkTimeList"
	TaskService_SaveTaskWorkTime_FullMethodName  = "/proto.TaskService/SaveTaskWorkTime"
	TaskService_SaveTaskFile_FullMethodName      = "/proto.TaskService/SaveTaskFile"
	TaskService_TaskSources_FullMethodName       = "/proto.TaskService/TaskSources"
	TaskService_CreateTaskComment_FullMethodName = "/proto.TaskService/CreateTaskComment"
	TaskService_CreateTaskList_FullMethodName    = "/proto.TaskService/CreateTaskList"
	TaskService_DeleteTaskStages_FullMethodName  = "/proto.TaskService/DeleteTaskStages"
	TaskService_EditTaskStages_FullMethodName    = "/proto.TaskService/EditTaskStages"
)

// TaskServiceClient is the client API for TaskService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskServiceClient interface {
	TaskStages(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskStagesResponse, error)
	MemberProjectList(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*MemberProjectResponse, error)
	TaskList(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskListResponse, error)
	SaveTask(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskMessage, error)
	TaskSort(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskSortResponse, error)
	MyTaskList(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*MyTaskListResponse, error)
	ReadTask(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskMessage, error)
	ListTaskMember(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskMemberList, error)
	TaskLog(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskLogList, error)
	TaskWorkTimeList(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskWorkTimeResponse, error)
	SaveTaskWorkTime(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*SaveTaskWorkTimeResponse, error)
	SaveTaskFile(ctx context.Context, in *TaskFileReqMessage, opts ...grpc.CallOption) (*TaskFileResponse, error)
	TaskSources(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskSourceResponse, error)
	CreateTaskComment(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error)
	CreateTaskList(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskStagesMessage, error)
	DeleteTaskStages(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*UpdateTaskStagesResponse, error)
	EditTaskStages(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*UpdateTaskStagesResponse, error)
}

type taskServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskServiceClient(cc grpc.ClientConnInterface) TaskServiceClient {
	return &taskServiceClient{cc}
}

func (c *taskServiceClient) TaskStages(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskStagesResponse, error) {
	out := new(TaskStagesResponse)
	err := c.cc.Invoke(ctx, TaskService_TaskStages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) MemberProjectList(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*MemberProjectResponse, error) {
	out := new(MemberProjectResponse)
	err := c.cc.Invoke(ctx, TaskService_MemberProjectList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) TaskList(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskListResponse, error) {
	out := new(TaskListResponse)
	err := c.cc.Invoke(ctx, TaskService_TaskList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) SaveTask(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskMessage, error) {
	out := new(TaskMessage)
	err := c.cc.Invoke(ctx, TaskService_SaveTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) TaskSort(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskSortResponse, error) {
	out := new(TaskSortResponse)
	err := c.cc.Invoke(ctx, TaskService_TaskSort_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) MyTaskList(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*MyTaskListResponse, error) {
	out := new(MyTaskListResponse)
	err := c.cc.Invoke(ctx, TaskService_MyTaskList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) ReadTask(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskMessage, error) {
	out := new(TaskMessage)
	err := c.cc.Invoke(ctx, TaskService_ReadTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) ListTaskMember(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskMemberList, error) {
	out := new(TaskMemberList)
	err := c.cc.Invoke(ctx, TaskService_ListTaskMember_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) TaskLog(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskLogList, error) {
	out := new(TaskLogList)
	err := c.cc.Invoke(ctx, TaskService_TaskLog_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) TaskWorkTimeList(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskWorkTimeResponse, error) {
	out := new(TaskWorkTimeResponse)
	err := c.cc.Invoke(ctx, TaskService_TaskWorkTimeList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) SaveTaskWorkTime(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*SaveTaskWorkTimeResponse, error) {
	out := new(SaveTaskWorkTimeResponse)
	err := c.cc.Invoke(ctx, TaskService_SaveTaskWorkTime_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) SaveTaskFile(ctx context.Context, in *TaskFileReqMessage, opts ...grpc.CallOption) (*TaskFileResponse, error) {
	out := new(TaskFileResponse)
	err := c.cc.Invoke(ctx, TaskService_SaveTaskFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) TaskSources(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskSourceResponse, error) {
	out := new(TaskSourceResponse)
	err := c.cc.Invoke(ctx, TaskService_TaskSources_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) CreateTaskComment(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error) {
	out := new(CreateCommentResponse)
	err := c.cc.Invoke(ctx, TaskService_CreateTaskComment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) CreateTaskList(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*TaskStagesMessage, error) {
	out := new(TaskStagesMessage)
	err := c.cc.Invoke(ctx, TaskService_CreateTaskList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) DeleteTaskStages(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*UpdateTaskStagesResponse, error) {
	out := new(UpdateTaskStagesResponse)
	err := c.cc.Invoke(ctx, TaskService_DeleteTaskStages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskServiceClient) EditTaskStages(ctx context.Context, in *TaskRequest, opts ...grpc.CallOption) (*UpdateTaskStagesResponse, error) {
	out := new(UpdateTaskStagesResponse)
	err := c.cc.Invoke(ctx, TaskService_EditTaskStages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskServiceServer is the server API for TaskService service.
// All implementations must embed UnimplementedTaskServiceServer
// for forward compatibility
type TaskServiceServer interface {
	TaskStages(context.Context, *TaskRequest) (*TaskStagesResponse, error)
	MemberProjectList(context.Context, *TaskRequest) (*MemberProjectResponse, error)
	TaskList(context.Context, *TaskRequest) (*TaskListResponse, error)
	SaveTask(context.Context, *TaskRequest) (*TaskMessage, error)
	TaskSort(context.Context, *TaskRequest) (*TaskSortResponse, error)
	MyTaskList(context.Context, *TaskRequest) (*MyTaskListResponse, error)
	ReadTask(context.Context, *TaskRequest) (*TaskMessage, error)
	ListTaskMember(context.Context, *TaskRequest) (*TaskMemberList, error)
	TaskLog(context.Context, *TaskRequest) (*TaskLogList, error)
	TaskWorkTimeList(context.Context, *TaskRequest) (*TaskWorkTimeResponse, error)
	SaveTaskWorkTime(context.Context, *TaskRequest) (*SaveTaskWorkTimeResponse, error)
	SaveTaskFile(context.Context, *TaskFileReqMessage) (*TaskFileResponse, error)
	TaskSources(context.Context, *TaskRequest) (*TaskSourceResponse, error)
	CreateTaskComment(context.Context, *TaskRequest) (*CreateCommentResponse, error)
	CreateTaskList(context.Context, *TaskRequest) (*TaskStagesMessage, error)
	DeleteTaskStages(context.Context, *TaskRequest) (*UpdateTaskStagesResponse, error)
	EditTaskStages(context.Context, *TaskRequest) (*UpdateTaskStagesResponse, error)
	mustEmbedUnimplementedTaskServiceServer()
}

// UnimplementedTaskServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTaskServiceServer struct {
}

func (UnimplementedTaskServiceServer) TaskStages(context.Context, *TaskRequest) (*TaskStagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskStages not implemented")
}
func (UnimplementedTaskServiceServer) MemberProjectList(context.Context, *TaskRequest) (*MemberProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MemberProjectList not implemented")
}
func (UnimplementedTaskServiceServer) TaskList(context.Context, *TaskRequest) (*TaskListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskList not implemented")
}
func (UnimplementedTaskServiceServer) SaveTask(context.Context, *TaskRequest) (*TaskMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveTask not implemented")
}
func (UnimplementedTaskServiceServer) TaskSort(context.Context, *TaskRequest) (*TaskSortResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskSort not implemented")
}
func (UnimplementedTaskServiceServer) MyTaskList(context.Context, *TaskRequest) (*MyTaskListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MyTaskList not implemented")
}
func (UnimplementedTaskServiceServer) ReadTask(context.Context, *TaskRequest) (*TaskMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadTask not implemented")
}
func (UnimplementedTaskServiceServer) ListTaskMember(context.Context, *TaskRequest) (*TaskMemberList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTaskMember not implemented")
}
func (UnimplementedTaskServiceServer) TaskLog(context.Context, *TaskRequest) (*TaskLogList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskLog not implemented")
}
func (UnimplementedTaskServiceServer) TaskWorkTimeList(context.Context, *TaskRequest) (*TaskWorkTimeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskWorkTimeList not implemented")
}
func (UnimplementedTaskServiceServer) SaveTaskWorkTime(context.Context, *TaskRequest) (*SaveTaskWorkTimeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveTaskWorkTime not implemented")
}
func (UnimplementedTaskServiceServer) SaveTaskFile(context.Context, *TaskFileReqMessage) (*TaskFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveTaskFile not implemented")
}
func (UnimplementedTaskServiceServer) TaskSources(context.Context, *TaskRequest) (*TaskSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskSources not implemented")
}
func (UnimplementedTaskServiceServer) CreateTaskComment(context.Context, *TaskRequest) (*CreateCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTaskComment not implemented")
}
func (UnimplementedTaskServiceServer) CreateTaskList(context.Context, *TaskRequest) (*TaskStagesMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTaskList not implemented")
}
func (UnimplementedTaskServiceServer) DeleteTaskStages(context.Context, *TaskRequest) (*UpdateTaskStagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTaskStages not implemented")
}
func (UnimplementedTaskServiceServer) EditTaskStages(context.Context, *TaskRequest) (*UpdateTaskStagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditTaskStages not implemented")
}
func (UnimplementedTaskServiceServer) mustEmbedUnimplementedTaskServiceServer() {}

// UnsafeTaskServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskServiceServer will
// result in compilation errors.
type UnsafeTaskServiceServer interface {
	mustEmbedUnimplementedTaskServiceServer()
}

func RegisterTaskServiceServer(s grpc.ServiceRegistrar, srv TaskServiceServer) {
	s.RegisterService(&TaskService_ServiceDesc, srv)
}

func _TaskService_TaskStages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).TaskStages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_TaskStages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).TaskStages(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_MemberProjectList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).MemberProjectList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_MemberProjectList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).MemberProjectList(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_TaskList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).TaskList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_TaskList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).TaskList(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_SaveTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).SaveTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_SaveTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).SaveTask(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_TaskSort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).TaskSort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_TaskSort_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).TaskSort(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_MyTaskList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).MyTaskList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_MyTaskList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).MyTaskList(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_ReadTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).ReadTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_ReadTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).ReadTask(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_ListTaskMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).ListTaskMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_ListTaskMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).ListTaskMember(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_TaskLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).TaskLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_TaskLog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).TaskLog(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_TaskWorkTimeList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).TaskWorkTimeList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_TaskWorkTimeList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).TaskWorkTimeList(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_SaveTaskWorkTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).SaveTaskWorkTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_SaveTaskWorkTime_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).SaveTaskWorkTime(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_SaveTaskFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskFileReqMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).SaveTaskFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_SaveTaskFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).SaveTaskFile(ctx, req.(*TaskFileReqMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_TaskSources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).TaskSources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_TaskSources_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).TaskSources(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_CreateTaskComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).CreateTaskComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_CreateTaskComment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).CreateTaskComment(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_CreateTaskList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).CreateTaskList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_CreateTaskList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).CreateTaskList(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_DeleteTaskStages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).DeleteTaskStages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_DeleteTaskStages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).DeleteTaskStages(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskService_EditTaskStages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskServiceServer).EditTaskStages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TaskService_EditTaskStages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskServiceServer).EditTaskStages(ctx, req.(*TaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TaskService_ServiceDesc is the grpc.ServiceDesc for TaskService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TaskService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TaskService",
	HandlerType: (*TaskServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TaskStages",
			Handler:    _TaskService_TaskStages_Handler,
		},
		{
			MethodName: "MemberProjectList",
			Handler:    _TaskService_MemberProjectList_Handler,
		},
		{
			MethodName: "TaskList",
			Handler:    _TaskService_TaskList_Handler,
		},
		{
			MethodName: "SaveTask",
			Handler:    _TaskService_SaveTask_Handler,
		},
		{
			MethodName: "TaskSort",
			Handler:    _TaskService_TaskSort_Handler,
		},
		{
			MethodName: "MyTaskList",
			Handler:    _TaskService_MyTaskList_Handler,
		},
		{
			MethodName: "ReadTask",
			Handler:    _TaskService_ReadTask_Handler,
		},
		{
			MethodName: "ListTaskMember",
			Handler:    _TaskService_ListTaskMember_Handler,
		},
		{
			MethodName: "TaskLog",
			Handler:    _TaskService_TaskLog_Handler,
		},
		{
			MethodName: "TaskWorkTimeList",
			Handler:    _TaskService_TaskWorkTimeList_Handler,
		},
		{
			MethodName: "SaveTaskWorkTime",
			Handler:    _TaskService_SaveTaskWorkTime_Handler,
		},
		{
			MethodName: "SaveTaskFile",
			Handler:    _TaskService_SaveTaskFile_Handler,
		},
		{
			MethodName: "TaskSources",
			Handler:    _TaskService_TaskSources_Handler,
		},
		{
			MethodName: "CreateTaskComment",
			Handler:    _TaskService_CreateTaskComment_Handler,
		},
		{
			MethodName: "CreateTaskList",
			Handler:    _TaskService_CreateTaskList_Handler,
		},
		{
			MethodName: "DeleteTaskStages",
			Handler:    _TaskService_DeleteTaskStages_Handler,
		},
		{
			MethodName: "EditTaskStages",
			Handler:    _TaskService_EditTaskStages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "task_service.proto",
}
