// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.1
// source: project_log.proto

package project_log

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

// ProjectLogServiceClient is the client API for ProjectLogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProjectLogServiceClient interface {
	PushLogs(ctx context.Context, in *PushLogsRequest, opts ...grpc.CallOption) (*PushLogsResponse, error)
}

type projectLogServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProjectLogServiceClient(cc grpc.ClientConnInterface) ProjectLogServiceClient {
	return &projectLogServiceClient{cc}
}

func (c *projectLogServiceClient) PushLogs(ctx context.Context, in *PushLogsRequest, opts ...grpc.CallOption) (*PushLogsResponse, error) {
	out := new(PushLogsResponse)
	err := c.cc.Invoke(ctx, "/project_log.ProjectLogService/PushLogs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProjectLogServiceServer is the server API for ProjectLogService service.
// All implementations must embed UnimplementedProjectLogServiceServer
// for forward compatibility
type ProjectLogServiceServer interface {
	PushLogs(context.Context, *PushLogsRequest) (*PushLogsResponse, error)
	mustEmbedUnimplementedProjectLogServiceServer()
}

// UnimplementedProjectLogServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProjectLogServiceServer struct {
}

func (UnimplementedProjectLogServiceServer) PushLogs(context.Context, *PushLogsRequest) (*PushLogsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushLogs not implemented")
}
func (UnimplementedProjectLogServiceServer) mustEmbedUnimplementedProjectLogServiceServer() {}

// UnsafeProjectLogServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProjectLogServiceServer will
// result in compilation errors.
type UnsafeProjectLogServiceServer interface {
	mustEmbedUnimplementedProjectLogServiceServer()
}

func RegisterProjectLogServiceServer(s grpc.ServiceRegistrar, srv ProjectLogServiceServer) {
	s.RegisterService(&ProjectLogService_ServiceDesc, srv)
}

func _ProjectLogService_PushLogs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushLogsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectLogServiceServer).PushLogs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project_log.ProjectLogService/PushLogs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectLogServiceServer).PushLogs(ctx, req.(*PushLogsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProjectLogService_ServiceDesc is the grpc.ServiceDesc for ProjectLogService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProjectLogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "project_log.ProjectLogService",
	HandlerType: (*ProjectLogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PushLogs",
			Handler:    _ProjectLogService_PushLogs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "project_log.proto",
}
