// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.27.0
// source: department_service.proto

package department_grpc

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
	DepartmentService_Save_FullMethodName           = "/department.service.v1.DepartmentService/Save"
	DepartmentService_List_FullMethodName           = "/department.service.v1.DepartmentService/List"
	DepartmentService_ReadDepartment_FullMethodName = "/department.service.v1.DepartmentService/ReadDepartment"
)

// DepartmentServiceClient is the client API for DepartmentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DepartmentServiceClient interface {
	Save(ctx context.Context, in *DepartmentReqMessage, opts ...grpc.CallOption) (*DepartmentMessage, error)
	List(ctx context.Context, in *DepartmentReqMessage, opts ...grpc.CallOption) (*ListDepartmentMessage, error)
	ReadDepartment(ctx context.Context, in *DepartmentReqMessage, opts ...grpc.CallOption) (*DepartmentMessage, error)
}

type departmentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDepartmentServiceClient(cc grpc.ClientConnInterface) DepartmentServiceClient {
	return &departmentServiceClient{cc}
}

func (c *departmentServiceClient) Save(ctx context.Context, in *DepartmentReqMessage, opts ...grpc.CallOption) (*DepartmentMessage, error) {
	out := new(DepartmentMessage)
	err := c.cc.Invoke(ctx, DepartmentService_Save_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *departmentServiceClient) List(ctx context.Context, in *DepartmentReqMessage, opts ...grpc.CallOption) (*ListDepartmentMessage, error) {
	out := new(ListDepartmentMessage)
	err := c.cc.Invoke(ctx, DepartmentService_List_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *departmentServiceClient) ReadDepartment(ctx context.Context, in *DepartmentReqMessage, opts ...grpc.CallOption) (*DepartmentMessage, error) {
	out := new(DepartmentMessage)
	err := c.cc.Invoke(ctx, DepartmentService_ReadDepartment_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DepartmentServiceServer is the server API for DepartmentService service.
// All implementations must embed UnimplementedDepartmentServiceServer
// for forward compatibility
type DepartmentServiceServer interface {
	Save(context.Context, *DepartmentReqMessage) (*DepartmentMessage, error)
	List(context.Context, *DepartmentReqMessage) (*ListDepartmentMessage, error)
	ReadDepartment(context.Context, *DepartmentReqMessage) (*DepartmentMessage, error)
	mustEmbedUnimplementedDepartmentServiceServer()
}

// UnimplementedDepartmentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDepartmentServiceServer struct {
}

func (UnimplementedDepartmentServiceServer) Save(context.Context, *DepartmentReqMessage) (*DepartmentMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Save not implemented")
}
func (UnimplementedDepartmentServiceServer) List(context.Context, *DepartmentReqMessage) (*ListDepartmentMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedDepartmentServiceServer) ReadDepartment(context.Context, *DepartmentReqMessage) (*DepartmentMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadDepartment not implemented")
}
func (UnimplementedDepartmentServiceServer) mustEmbedUnimplementedDepartmentServiceServer() {}

// UnsafeDepartmentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DepartmentServiceServer will
// result in compilation errors.
type UnsafeDepartmentServiceServer interface {
	mustEmbedUnimplementedDepartmentServiceServer()
}

func RegisterDepartmentServiceServer(s grpc.ServiceRegistrar, srv DepartmentServiceServer) {
	s.RegisterService(&DepartmentService_ServiceDesc, srv)
}

func _DepartmentService_Save_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepartmentReqMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepartmentServiceServer).Save(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DepartmentService_Save_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepartmentServiceServer).Save(ctx, req.(*DepartmentReqMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _DepartmentService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepartmentReqMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepartmentServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DepartmentService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepartmentServiceServer).List(ctx, req.(*DepartmentReqMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _DepartmentService_ReadDepartment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepartmentReqMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DepartmentServiceServer).ReadDepartment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DepartmentService_ReadDepartment_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DepartmentServiceServer).ReadDepartment(ctx, req.(*DepartmentReqMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// DepartmentService_ServiceDesc is the grpc.ServiceDesc for DepartmentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DepartmentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "department.service.v1.DepartmentService",
	HandlerType: (*DepartmentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Save",
			Handler:    _DepartmentService_Save_Handler,
		},
		{
			MethodName: "List",
			Handler:    _DepartmentService_List_Handler,
		},
		{
			MethodName: "ReadDepartment",
			Handler:    _DepartmentService_ReadDepartment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "department_service.proto",
}
