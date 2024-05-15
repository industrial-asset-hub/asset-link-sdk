// ------------------------------------------------------------------
// Common Definition of Device Identifiers
// ------------------------------------------------------------------

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.2
// source: common_identifiers.proto

package iah_discovery

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
	IdentifiersApi_GetIdentifiers_FullMethodName        = "/siemens.common.identifiers.v1.IdentifiersApi/GetIdentifiers"
	IdentifiersApi_GetSupportedSemantics_FullMethodName = "/siemens.common.identifiers.v1.IdentifiersApi/GetSupportedSemantics"
)

// IdentifiersApiClient is the client API for IdentifiersApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IdentifiersApiClient interface {
	GetIdentifiers(ctx context.Context, in *GetIdentifiersRequest, opts ...grpc.CallOption) (*GetIdentifiersResponse, error)
	GetSupportedSemantics(ctx context.Context, in *GetSupportedSemanticsRequest, opts ...grpc.CallOption) (*GetSupportedSemanticsResponse, error)
}

type identifiersApiClient struct {
	cc grpc.ClientConnInterface
}

func NewIdentifiersApiClient(cc grpc.ClientConnInterface) IdentifiersApiClient {
	return &identifiersApiClient{cc}
}

func (c *identifiersApiClient) GetIdentifiers(ctx context.Context, in *GetIdentifiersRequest, opts ...grpc.CallOption) (*GetIdentifiersResponse, error) {
	out := new(GetIdentifiersResponse)
	err := c.cc.Invoke(ctx, IdentifiersApi_GetIdentifiers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identifiersApiClient) GetSupportedSemantics(ctx context.Context, in *GetSupportedSemanticsRequest, opts ...grpc.CallOption) (*GetSupportedSemanticsResponse, error) {
	out := new(GetSupportedSemanticsResponse)
	err := c.cc.Invoke(ctx, IdentifiersApi_GetSupportedSemantics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IdentifiersApiServer is the server API for IdentifiersApi service.
// All implementations must embed UnimplementedIdentifiersApiServer
// for forward compatibility
type IdentifiersApiServer interface {
	GetIdentifiers(context.Context, *GetIdentifiersRequest) (*GetIdentifiersResponse, error)
	GetSupportedSemantics(context.Context, *GetSupportedSemanticsRequest) (*GetSupportedSemanticsResponse, error)
	mustEmbedUnimplementedIdentifiersApiServer()
}

// UnimplementedIdentifiersApiServer must be embedded to have forward compatible implementations.
type UnimplementedIdentifiersApiServer struct {
}

func (UnimplementedIdentifiersApiServer) GetIdentifiers(context.Context, *GetIdentifiersRequest) (*GetIdentifiersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIdentifiers not implemented")
}
func (UnimplementedIdentifiersApiServer) GetSupportedSemantics(context.Context, *GetSupportedSemanticsRequest) (*GetSupportedSemanticsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSupportedSemantics not implemented")
}
func (UnimplementedIdentifiersApiServer) mustEmbedUnimplementedIdentifiersApiServer() {}

// UnsafeIdentifiersApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IdentifiersApiServer will
// result in compilation errors.
type UnsafeIdentifiersApiServer interface {
	mustEmbedUnimplementedIdentifiersApiServer()
}

func RegisterIdentifiersApiServer(s grpc.ServiceRegistrar, srv IdentifiersApiServer) {
	s.RegisterService(&IdentifiersApi_ServiceDesc, srv)
}

func _IdentifiersApi_GetIdentifiers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIdentifiersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentifiersApiServer).GetIdentifiers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IdentifiersApi_GetIdentifiers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentifiersApiServer).GetIdentifiers(ctx, req.(*GetIdentifiersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentifiersApi_GetSupportedSemantics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSupportedSemanticsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentifiersApiServer).GetSupportedSemantics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IdentifiersApi_GetSupportedSemantics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentifiersApiServer).GetSupportedSemantics(ctx, req.(*GetSupportedSemanticsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IdentifiersApi_ServiceDesc is the grpc.ServiceDesc for IdentifiersApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IdentifiersApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "siemens.common.identifiers.v1.IdentifiersApi",
	HandlerType: (*IdentifiersApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetIdentifiers",
			Handler:    _IdentifiersApi_GetIdentifiers_Handler,
		},
		{
			MethodName: "GetSupportedSemantics",
			Handler:    _IdentifiersApi_GetSupportedSemantics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "common_identifiers.proto",
}
