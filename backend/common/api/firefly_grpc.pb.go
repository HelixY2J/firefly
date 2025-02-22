// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: api/firefly.proto

package api

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
	FireflyService_RegisterNode_FullMethodName    = "/api.FireflyService/RegisterNode"
	FireflyService_Heartbeat_FullMethodName       = "/api.FireflyService/Heartbeat"
	FireflyService_SyncLibrary_FullMethodName     = "/api.FireflyService/SyncLibrary"
	FireflyService_RequestPlayback_FullMethodName = "/api.FireflyService/RequestPlayback"
	FireflyService_SyncPlayback_FullMethodName    = "/api.FireflyService/SyncPlayback"
)

// FireflyServiceClient is the client API for FireflyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FireflyServiceClient interface {
	RegisterNode(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error)
	SyncLibrary(ctx context.Context, in *SyncLibraryRequest, opts ...grpc.CallOption) (*SyncLibraryResponse, error)
	RequestPlayback(ctx context.Context, in *PlaybackRequest, opts ...grpc.CallOption) (*PlaybackResponse, error)
	SyncPlayback(ctx context.Context, in *SyncPlaybackCommand, opts ...grpc.CallOption) (*SyncPlaybackResponse, error)
}

type fireflyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFireflyServiceClient(cc grpc.ClientConnInterface) FireflyServiceClient {
	return &fireflyServiceClient{cc}
}

func (c *fireflyServiceClient) RegisterNode(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, FireflyService_RegisterNode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fireflyServiceClient) Heartbeat(ctx context.Context, in *HeartbeatRequest, opts ...grpc.CallOption) (*HeartbeatResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HeartbeatResponse)
	err := c.cc.Invoke(ctx, FireflyService_Heartbeat_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fireflyServiceClient) SyncLibrary(ctx context.Context, in *SyncLibraryRequest, opts ...grpc.CallOption) (*SyncLibraryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SyncLibraryResponse)
	err := c.cc.Invoke(ctx, FireflyService_SyncLibrary_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fireflyServiceClient) RequestPlayback(ctx context.Context, in *PlaybackRequest, opts ...grpc.CallOption) (*PlaybackResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PlaybackResponse)
	err := c.cc.Invoke(ctx, FireflyService_RequestPlayback_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fireflyServiceClient) SyncPlayback(ctx context.Context, in *SyncPlaybackCommand, opts ...grpc.CallOption) (*SyncPlaybackResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SyncPlaybackResponse)
	err := c.cc.Invoke(ctx, FireflyService_SyncPlayback_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FireflyServiceServer is the server API for FireflyService service.
// All implementations must embed UnimplementedFireflyServiceServer
// for forward compatibility.
type FireflyServiceServer interface {
	RegisterNode(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error)
	SyncLibrary(context.Context, *SyncLibraryRequest) (*SyncLibraryResponse, error)
	RequestPlayback(context.Context, *PlaybackRequest) (*PlaybackResponse, error)
	SyncPlayback(context.Context, *SyncPlaybackCommand) (*SyncPlaybackResponse, error)
	mustEmbedUnimplementedFireflyServiceServer()
}

// UnimplementedFireflyServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFireflyServiceServer struct{}

func (UnimplementedFireflyServiceServer) RegisterNode(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterNode not implemented")
}
func (UnimplementedFireflyServiceServer) Heartbeat(context.Context, *HeartbeatRequest) (*HeartbeatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Heartbeat not implemented")
}
func (UnimplementedFireflyServiceServer) SyncLibrary(context.Context, *SyncLibraryRequest) (*SyncLibraryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncLibrary not implemented")
}
func (UnimplementedFireflyServiceServer) RequestPlayback(context.Context, *PlaybackRequest) (*PlaybackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestPlayback not implemented")
}
func (UnimplementedFireflyServiceServer) SyncPlayback(context.Context, *SyncPlaybackCommand) (*SyncPlaybackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncPlayback not implemented")
}
func (UnimplementedFireflyServiceServer) mustEmbedUnimplementedFireflyServiceServer() {}
func (UnimplementedFireflyServiceServer) testEmbeddedByValue()                        {}

// UnsafeFireflyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FireflyServiceServer will
// result in compilation errors.
type UnsafeFireflyServiceServer interface {
	mustEmbedUnimplementedFireflyServiceServer()
}

func RegisterFireflyServiceServer(s grpc.ServiceRegistrar, srv FireflyServiceServer) {
	// If the following call pancis, it indicates UnimplementedFireflyServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FireflyService_ServiceDesc, srv)
}

func _FireflyService_RegisterNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FireflyServiceServer).RegisterNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FireflyService_RegisterNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FireflyServiceServer).RegisterNode(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FireflyService_Heartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartbeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FireflyServiceServer).Heartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FireflyService_Heartbeat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FireflyServiceServer).Heartbeat(ctx, req.(*HeartbeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FireflyService_SyncLibrary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncLibraryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FireflyServiceServer).SyncLibrary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FireflyService_SyncLibrary_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FireflyServiceServer).SyncLibrary(ctx, req.(*SyncLibraryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FireflyService_RequestPlayback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaybackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FireflyServiceServer).RequestPlayback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FireflyService_RequestPlayback_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FireflyServiceServer).RequestPlayback(ctx, req.(*PlaybackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FireflyService_SyncPlayback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncPlaybackCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FireflyServiceServer).SyncPlayback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FireflyService_SyncPlayback_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FireflyServiceServer).SyncPlayback(ctx, req.(*SyncPlaybackCommand))
	}
	return interceptor(ctx, in, info, handler)
}

// FireflyService_ServiceDesc is the grpc.ServiceDesc for FireflyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FireflyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.FireflyService",
	HandlerType: (*FireflyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterNode",
			Handler:    _FireflyService_RegisterNode_Handler,
		},
		{
			MethodName: "Heartbeat",
			Handler:    _FireflyService_Heartbeat_Handler,
		},
		{
			MethodName: "SyncLibrary",
			Handler:    _FireflyService_SyncLibrary_Handler,
		},
		{
			MethodName: "RequestPlayback",
			Handler:    _FireflyService_RequestPlayback_Handler,
		},
		{
			MethodName: "SyncPlayback",
			Handler:    _FireflyService_SyncPlayback_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/firefly.proto",
}
