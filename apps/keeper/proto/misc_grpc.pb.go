// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: proto/misc.proto

package proto

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

// KeeperAPIServiceClient is the client API for KeeperAPIService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeeperAPIServiceClient interface {
	HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type keeperAPIServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKeeperAPIServiceClient(cc grpc.ClientConnInterface) KeeperAPIServiceClient {
	return &keeperAPIServiceClient{cc}
}

func (c *keeperAPIServiceClient) HealthCheck(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, "/keeper_api.KeeperAPIService/HealthCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeeperAPIServiceServer is the server API for KeeperAPIService service.
// All implementations must embed UnimplementedKeeperAPIServiceServer
// for forward compatibility
type KeeperAPIServiceServer interface {
	HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
	mustEmbedUnimplementedKeeperAPIServiceServer()
}

// UnimplementedKeeperAPIServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKeeperAPIServiceServer struct {
}

func (UnimplementedKeeperAPIServiceServer) HealthCheck(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedKeeperAPIServiceServer) mustEmbedUnimplementedKeeperAPIServiceServer() {}

// UnsafeKeeperAPIServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeeperAPIServiceServer will
// result in compilation errors.
type UnsafeKeeperAPIServiceServer interface {
	mustEmbedUnimplementedKeeperAPIServiceServer()
}

func RegisterKeeperAPIServiceServer(s grpc.ServiceRegistrar, srv KeeperAPIServiceServer) {
	s.RegisterService(&KeeperAPIService_ServiceDesc, srv)
}

func _KeeperAPIService_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperAPIServiceServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/keeper_api.KeeperAPIService/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperAPIServiceServer).HealthCheck(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KeeperAPIService_ServiceDesc is the grpc.ServiceDesc for KeeperAPIService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeeperAPIService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "keeper_api.KeeperAPIService",
	HandlerType: (*KeeperAPIServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _KeeperAPIService_HealthCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/misc.proto",
}
