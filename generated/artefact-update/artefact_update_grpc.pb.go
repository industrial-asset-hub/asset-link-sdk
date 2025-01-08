// Artefact update interface
// This is the interface for pushing and pulling
// artefacts to and from drivers.
// The driver is responsible for the actual
// transfer of the artefact to the target device.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: artefact_update.proto

package artefact_update

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
	ArtefactUpdateApi_PushArtefact_FullMethodName = "/factory_x.artefact_update.v1.ArtefactUpdateApi/PushArtefact"
	ArtefactUpdateApi_PullArtefact_FullMethodName = "/factory_x.artefact_update.v1.ArtefactUpdateApi/PullArtefact"
)

// ArtefactUpdateApiClient is the client API for ArtefactUpdateApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ArtefactUpdateApiClient interface {
	// Push an artifact to a driver
	PushArtefact(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ArtefactChunk, ArtefactUpdateStatus], error)
	// Load an artifact from a driver
	PullArtefact(ctx context.Context, in *ArtefactType, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ArtefactChunk], error)
}

type artefactUpdateApiClient struct {
	cc grpc.ClientConnInterface
}

func NewArtefactUpdateApiClient(cc grpc.ClientConnInterface) ArtefactUpdateApiClient {
	return &artefactUpdateApiClient{cc}
}

func (c *artefactUpdateApiClient) PushArtefact(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ArtefactChunk, ArtefactUpdateStatus], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ArtefactUpdateApi_ServiceDesc.Streams[0], ArtefactUpdateApi_PushArtefact_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ArtefactChunk, ArtefactUpdateStatus]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ArtefactUpdateApi_PushArtefactClient = grpc.BidiStreamingClient[ArtefactChunk, ArtefactUpdateStatus]

func (c *artefactUpdateApiClient) PullArtefact(ctx context.Context, in *ArtefactType, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ArtefactChunk], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &ArtefactUpdateApi_ServiceDesc.Streams[1], ArtefactUpdateApi_PullArtefact_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ArtefactType, ArtefactChunk]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ArtefactUpdateApi_PullArtefactClient = grpc.ServerStreamingClient[ArtefactChunk]

// ArtefactUpdateApiServer is the server API for ArtefactUpdateApi service.
// All implementations must embed UnimplementedArtefactUpdateApiServer
// for forward compatibility.
type ArtefactUpdateApiServer interface {
	// Push an artifact to a driver
	PushArtefact(grpc.BidiStreamingServer[ArtefactChunk, ArtefactUpdateStatus]) error
	// Load an artifact from a driver
	PullArtefact(*ArtefactType, grpc.ServerStreamingServer[ArtefactChunk]) error
	mustEmbedUnimplementedArtefactUpdateApiServer()
}

// UnimplementedArtefactUpdateApiServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedArtefactUpdateApiServer struct{}

func (UnimplementedArtefactUpdateApiServer) PushArtefact(grpc.BidiStreamingServer[ArtefactChunk, ArtefactUpdateStatus]) error {
	return status.Errorf(codes.Unimplemented, "method PushArtefact not implemented")
}
func (UnimplementedArtefactUpdateApiServer) PullArtefact(*ArtefactType, grpc.ServerStreamingServer[ArtefactChunk]) error {
	return status.Errorf(codes.Unimplemented, "method PullArtefact not implemented")
}
func (UnimplementedArtefactUpdateApiServer) mustEmbedUnimplementedArtefactUpdateApiServer() {}
func (UnimplementedArtefactUpdateApiServer) testEmbeddedByValue()                           {}

// UnsafeArtefactUpdateApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ArtefactUpdateApiServer will
// result in compilation errors.
type UnsafeArtefactUpdateApiServer interface {
	mustEmbedUnimplementedArtefactUpdateApiServer()
}

func RegisterArtefactUpdateApiServer(s grpc.ServiceRegistrar, srv ArtefactUpdateApiServer) {
	// If the following call pancis, it indicates UnimplementedArtefactUpdateApiServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ArtefactUpdateApi_ServiceDesc, srv)
}

func _ArtefactUpdateApi_PushArtefact_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ArtefactUpdateApiServer).PushArtefact(&grpc.GenericServerStream[ArtefactChunk, ArtefactUpdateStatus]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ArtefactUpdateApi_PushArtefactServer = grpc.BidiStreamingServer[ArtefactChunk, ArtefactUpdateStatus]

func _ArtefactUpdateApi_PullArtefact_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ArtefactType)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ArtefactUpdateApiServer).PullArtefact(m, &grpc.GenericServerStream[ArtefactType, ArtefactChunk]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type ArtefactUpdateApi_PullArtefactServer = grpc.ServerStreamingServer[ArtefactChunk]

// ArtefactUpdateApi_ServiceDesc is the grpc.ServiceDesc for ArtefactUpdateApi service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ArtefactUpdateApi_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "factory_x.artefact_update.v1.ArtefactUpdateApi",
	HandlerType: (*ArtefactUpdateApiServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PushArtefact",
			Handler:       _ArtefactUpdateApi_PushArtefact_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "PullArtefact",
			Handler:       _ArtefactUpdateApi_PullArtefact_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "artefact_update.proto",
}
