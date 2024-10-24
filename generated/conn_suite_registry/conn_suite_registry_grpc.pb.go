// ------------------------------------------------------------------
// Connectivity Suite Registry
// ------------------------------------------------------------------
//
// Naming convention according:
// https://cloud.google.com/apis/design/naming_convention
//
// ------------------------------------------------------------------

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.28.2
// source: conn_suite_registry.proto

package conn_suite_registry

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
	RegistryApi_RegisterService_FullMethodName         = "/siemens.connectivitysuite.registry.v1.RegistryApi/RegisterService"
	RegistryApi_UnregisterService_FullMethodName       = "/siemens.connectivitysuite.registry.v1.RegistryApi/UnregisterService"
	RegistryApi_QueryRegisteredServices_FullMethodName = "/siemens.connectivitysuite.registry.v1.RegistryApi/QueryRegisteredServices"
)

// RegistryApiClient is the client API for RegistryApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegistryApiClient interface {
	// Possible return values
	// - OK
	RegisterService(ctx context.Context, in *RegisterServiceRequest, opts ...grpc.CallOption) (*RegisterServiceResponse, error)
	// Possible return values
	// - OK
	// - NOT_FOUND - service with the provided key was not registered
	UnregisterService(ctx context.Context, in *UnregisterServiceRequest, opts ...grpc.CallOption) (*UnregisterServiceResponse, error)
	// Possible return values
	// - OK
	QueryRegisteredServices(ctx context.Context, in *QueryRegisteredServicesRequest, opts ...grpc.CallOption) (*QueryRegisteredServicesResponse, error)
}

type registryApiClient struct {
	cc grpc.ClientConnInterface
}

func NewRegistryApiClient(cc grpc.ClientConnInterface) RegistryApiClient {
	return &registryApiClient{cc}
}

func (c *registryApiClient) RegisterService(ctx context.Context, in *RegisterServiceRequest, opts ...grpc.CallOption) (*RegisterServiceResponse, error) {
	out := new(RegisterServiceResponse)
	err := c.cc.Invoke(ctx, RegistryApi_RegisterService_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryApiClient) UnregisterService(ctx context.Context, in *UnregisterServiceRequest, opts ...grpc.CallOption) (*UnregisterServiceResponse, error) {
	out := new(UnregisterServiceResponse)
	err := c.cc.Invoke(ctx, RegistryApi_UnregisterService_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryApiClient) QueryRegisteredServices(ctx context.Context, in *QueryRegisteredServicesRequest, opts ...grpc.CallOption) (*QueryRegisteredServicesResponse, error) {
	out := new(QueryRegisteredServicesResponse)
	err := c.cc.Invoke(ctx, RegistryApi_QueryRegisteredServices_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegistryApiServer is the server API for RegistryApi service.
// All implementations must embed UnimplementedRegistryApiServer
// for forward compatibility
type RegistryApiServer interface {
	// Possible return values
	// - OK
	RegisterService(context.Context, *RegisterServiceRequest) (*RegisterServiceResponse, error)
	// Possible return values
	// - OK
	// - NOT_FOUND - service with the provided key was not registered
	UnregisterService(context.Context, *UnregisterServiceRequest) (*UnregisterServiceResponse, error)
	// Possible return values
	// - OK
	QueryRegisteredServices(context.Context, *QueryRegisteredServicesRequest) (*QueryRegisteredServicesResponse, error)
	mustEmbedUnimplementedRegistryApiServer()
}

// UnimplementedRegistryApiServer must be embedded to have forward compatible implementations.
type UnimplementedRegistryApiServer struct {
}

func (UnimplementedRegistryApiServer) RegisterService(context.Context, *RegisterServiceRequest) (*RegisterServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterService not implemented")
}
func (UnimplementedRegistryApiServer) UnregisterService(context.Context, *UnregisterServiceRequest) (*UnregisterServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnregisterService not implemented")
}
func (UnimplementedRegistryApiServer) QueryRegisteredServices(context.Context, *QueryRegisteredServicesRequest) (*QueryRegisteredServicesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryRegisteredServices not implemented")
}
func (UnimplementedRegistryApiServer) mustEmbedUnimplementedRegistryApiServer() {}

// UnsafeRegistryApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegistryApiServer will
// result in compilation errors.
type UnsafeRegistryApiServer interface {
	mustEmbedUnimplementedRegistryApiServer()
}

func RegisterRegistryApiServer(s grpc.ServiceRegistrar, srv RegistryApiServer) {
	s.RegisterService(&RegistryApi_ServiceDesc, srv)
}

func _RegistryApi_RegisterService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryApiServer).RegisterService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryApi_RegisterService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryApiServer).RegisterService(ctx, req.(*RegisterServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryApi_UnregisterService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnregisterServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryApiServer).UnregisterService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryApi_UnregisterService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryApiServer).UnregisterService(ctx, req.(*UnregisterServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RegistryApi_QueryRegisteredServices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRegisteredServicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryApiServer).QueryRegisteredServices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RegistryApi_QueryRegisteredServices_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryApiServer).QueryRegisteredServices(ctx, req.(*QueryRegisteredServicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegistryApi_ServiceDesc is the grpc.ServiceDesc for RegistryApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RegistryApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "siemens.connectivitysuite.registry.v1.RegistryApi",
	HandlerType: (*RegistryApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterService",
			Handler:    _RegistryApi_RegisterService_Handler,
		},
		{
			MethodName: "UnregisterService",
			Handler:    _RegistryApi_UnregisterService_Handler,
		},
		{
			MethodName: "QueryRegisteredServices",
			Handler:    _RegistryApi_QueryRegisteredServices_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "conn_suite_registry.proto",
}
