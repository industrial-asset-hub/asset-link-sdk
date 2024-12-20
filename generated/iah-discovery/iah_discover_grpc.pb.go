// Device Discover Interface

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.2
// source: iah_discover.proto

package iah_discovery

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DeviceDiscoverApi_GetFilterTypes_FullMethodName   = "/siemens.industrialassethub.discover.v1.DeviceDiscoverApi/GetFilterTypes"
	DeviceDiscoverApi_GetFilterOptions_FullMethodName = "/siemens.industrialassethub.discover.v1.DeviceDiscoverApi/GetFilterOptions"
	DeviceDiscoverApi_DiscoverDevices_FullMethodName  = "/siemens.industrialassethub.discover.v1.DeviceDiscoverApi/DiscoverDevices"
)

// DeviceDiscoverApiClient is the client API for DeviceDiscoverApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeviceDiscoverApiClient interface {
	// Get the list of supported Filter Types
	GetFilterTypes(ctx context.Context, in *FilterTypesRequest, opts ...grpc.CallOption) (*FilterTypesResponse, error)
	// Get the list of supported Filter Options
	GetFilterOptions(ctx context.Context, in *FilterOptionsRequest, opts ...grpc.CallOption) (*FilterOptionsResponse, error)
	// Start a device discovery with given filters and options.
	// Returns the discovered devices.
	DiscoverDevices(ctx context.Context, in *DiscoverRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DiscoverResponse], error)
}

type deviceDiscoverApiClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceDiscoverApiClient(cc grpc.ClientConnInterface) DeviceDiscoverApiClient {
	return &deviceDiscoverApiClient{cc}
}

func (c *deviceDiscoverApiClient) GetFilterTypes(ctx context.Context, in *FilterTypesRequest, opts ...grpc.CallOption) (*FilterTypesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FilterTypesResponse)
	err := c.cc.Invoke(ctx, DeviceDiscoverApi_GetFilterTypes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceDiscoverApiClient) GetFilterOptions(ctx context.Context, in *FilterOptionsRequest, opts ...grpc.CallOption) (*FilterOptionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FilterOptionsResponse)
	err := c.cc.Invoke(ctx, DeviceDiscoverApi_GetFilterOptions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceDiscoverApiClient) DiscoverDevices(ctx context.Context, in *DiscoverRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DiscoverResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DeviceDiscoverApi_ServiceDesc.Streams[0], DeviceDiscoverApi_DiscoverDevices_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DiscoverRequest, DiscoverResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DeviceDiscoverApi_DiscoverDevicesClient = grpc.ServerStreamingClient[DiscoverResponse]

// DeviceDiscoverApiServer is the server API for DeviceDiscoverApi service.
// All implementations must embed UnimplementedDeviceDiscoverApiServer
// for forward compatibility.
type DeviceDiscoverApiServer interface {
	// Get the list of supported Filter Types
	GetFilterTypes(context.Context, *FilterTypesRequest) (*FilterTypesResponse, error)
	// Get the list of supported Filter Options
	GetFilterOptions(context.Context, *FilterOptionsRequest) (*FilterOptionsResponse, error)
	// Start a device discovery with given filters and options.
	// Returns the discovered devices.
	DiscoverDevices(*DiscoverRequest, grpc.ServerStreamingServer[DiscoverResponse]) error
	mustEmbedUnimplementedDeviceDiscoverApiServer()
}

// UnimplementedDeviceDiscoverApiServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDeviceDiscoverApiServer struct{}

func (UnimplementedDeviceDiscoverApiServer) GetFilterTypes(context.Context, *FilterTypesRequest) (*FilterTypesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFilterTypes not implemented")
}
func (UnimplementedDeviceDiscoverApiServer) GetFilterOptions(context.Context, *FilterOptionsRequest) (*FilterOptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFilterOptions not implemented")
}
func (UnimplementedDeviceDiscoverApiServer) DiscoverDevices(*DiscoverRequest, grpc.ServerStreamingServer[DiscoverResponse]) error {
	return status.Errorf(codes.Unimplemented, "method DiscoverDevices not implemented")
}
func (UnimplementedDeviceDiscoverApiServer) mustEmbedUnimplementedDeviceDiscoverApiServer() {}
func (UnimplementedDeviceDiscoverApiServer) testEmbeddedByValue()                           {}

// UnsafeDeviceDiscoverApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceDiscoverApiServer will
// result in compilation errors.
type UnsafeDeviceDiscoverApiServer interface {
	mustEmbedUnimplementedDeviceDiscoverApiServer()
}

func RegisterDeviceDiscoverApiServer(s grpc.ServiceRegistrar, srv DeviceDiscoverApiServer) {
	// If the following call pancis, it indicates UnimplementedDeviceDiscoverApiServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DeviceDiscoverApi_ServiceDesc, srv)
}

func _DeviceDiscoverApi_GetFilterTypes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FilterTypesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceDiscoverApiServer).GetFilterTypes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceDiscoverApi_GetFilterTypes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceDiscoverApiServer).GetFilterTypes(ctx, req.(*FilterTypesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceDiscoverApi_GetFilterOptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FilterOptionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceDiscoverApiServer).GetFilterOptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeviceDiscoverApi_GetFilterOptions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceDiscoverApiServer).GetFilterOptions(ctx, req.(*FilterOptionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceDiscoverApi_DiscoverDevices_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DiscoverRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DeviceDiscoverApiServer).DiscoverDevices(m, &grpc.GenericServerStream[DiscoverRequest, DiscoverResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DeviceDiscoverApi_DiscoverDevicesServer = grpc.ServerStreamingServer[DiscoverResponse]

// DeviceDiscoverApi_ServiceDesc is the grpc.ServiceDesc for DeviceDiscoverApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeviceDiscoverApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "siemens.industrialassethub.discover.v1.DeviceDiscoverApi",
	HandlerType: (*DeviceDiscoverApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFilterTypes",
			Handler:    _DeviceDiscoverApi_GetFilterTypes_Handler,
		},
		{
			MethodName: "GetFilterOptions",
			Handler:    _DeviceDiscoverApi_GetFilterOptions_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DiscoverDevices",
			Handler:       _DeviceDiscoverApi_DiscoverDevices_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "iah_discover.proto",
}
