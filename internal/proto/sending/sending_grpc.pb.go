// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: sending.proto

package sending

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
	Sender_SendEmailVerification_FullMethodName = "/sending.Sender/SendEmailVerification"
	Sender_SendPasswordRecovery_FullMethodName  = "/sending.Sender/SendPasswordRecovery"
	Sender_SendUsernameRecovery_FullMethodName  = "/sending.Sender/SendUsernameRecovery"
)

// SenderClient is the client API for Sender service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SenderClient interface {
	SendEmailVerification(ctx context.Context, in *SendEmailVerificationRequest, opts ...grpc.CallOption) (*SendEmailReply, error)
	SendPasswordRecovery(ctx context.Context, in *SendPasswordRecoveryRequest, opts ...grpc.CallOption) (*SendEmailReply, error)
	SendUsernameRecovery(ctx context.Context, in *SendUsernameRecoveryRequest, opts ...grpc.CallOption) (*SendEmailReply, error)
}

type senderClient struct {
	cc grpc.ClientConnInterface
}

func NewSenderClient(cc grpc.ClientConnInterface) SenderClient {
	return &senderClient{cc}
}

func (c *senderClient) SendEmailVerification(ctx context.Context, in *SendEmailVerificationRequest, opts ...grpc.CallOption) (*SendEmailReply, error) {
	out := new(SendEmailReply)
	err := c.cc.Invoke(ctx, Sender_SendEmailVerification_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *senderClient) SendPasswordRecovery(ctx context.Context, in *SendPasswordRecoveryRequest, opts ...grpc.CallOption) (*SendEmailReply, error) {
	out := new(SendEmailReply)
	err := c.cc.Invoke(ctx, Sender_SendPasswordRecovery_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *senderClient) SendUsernameRecovery(ctx context.Context, in *SendUsernameRecoveryRequest, opts ...grpc.CallOption) (*SendEmailReply, error) {
	out := new(SendEmailReply)
	err := c.cc.Invoke(ctx, Sender_SendUsernameRecovery_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SenderServer is the server API for Sender service.
// All implementations must embed UnimplementedSenderServer
// for forward compatibility
type SenderServer interface {
	SendEmailVerification(context.Context, *SendEmailVerificationRequest) (*SendEmailReply, error)
	SendPasswordRecovery(context.Context, *SendPasswordRecoveryRequest) (*SendEmailReply, error)
	SendUsernameRecovery(context.Context, *SendUsernameRecoveryRequest) (*SendEmailReply, error)
	mustEmbedUnimplementedSenderServer()
}

// UnimplementedSenderServer must be embedded to have forward compatible implementations.
type UnimplementedSenderServer struct {
}

func (UnimplementedSenderServer) SendEmailVerification(context.Context, *SendEmailVerificationRequest) (*SendEmailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmailVerification not implemented")
}
func (UnimplementedSenderServer) SendPasswordRecovery(context.Context, *SendPasswordRecoveryRequest) (*SendEmailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPasswordRecovery not implemented")
}
func (UnimplementedSenderServer) SendUsernameRecovery(context.Context, *SendUsernameRecoveryRequest) (*SendEmailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendUsernameRecovery not implemented")
}
func (UnimplementedSenderServer) mustEmbedUnimplementedSenderServer() {}

// UnsafeSenderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SenderServer will
// result in compilation errors.
type UnsafeSenderServer interface {
	mustEmbedUnimplementedSenderServer()
}

func RegisterSenderServer(s grpc.ServiceRegistrar, srv SenderServer) {
	s.RegisterService(&Sender_ServiceDesc, srv)
}

func _Sender_SendEmailVerification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendEmailVerificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SenderServer).SendEmailVerification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sender_SendEmailVerification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SenderServer).SendEmailVerification(ctx, req.(*SendEmailVerificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sender_SendPasswordRecovery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendPasswordRecoveryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SenderServer).SendPasswordRecovery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sender_SendPasswordRecovery_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SenderServer).SendPasswordRecovery(ctx, req.(*SendPasswordRecoveryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sender_SendUsernameRecovery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendUsernameRecoveryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SenderServer).SendUsernameRecovery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Sender_SendUsernameRecovery_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SenderServer).SendUsernameRecovery(ctx, req.(*SendUsernameRecoveryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Sender_ServiceDesc is the grpc.ServiceDesc for Sender service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sender_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sending.Sender",
	HandlerType: (*SenderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEmailVerification",
			Handler:    _Sender_SendEmailVerification_Handler,
		},
		{
			MethodName: "SendPasswordRecovery",
			Handler:    _Sender_SendPasswordRecovery_Handler,
		},
		{
			MethodName: "SendUsernameRecovery",
			Handler:    _Sender_SendUsernameRecovery_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sending.proto",
}