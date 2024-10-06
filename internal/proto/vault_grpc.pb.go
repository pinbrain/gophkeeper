// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: internal/proto/vault.proto

package proto

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
	VaultService_Test_FullMethodName = "/VaultService/Test"
)

// VaultServiceClient is the client API for VaultService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VaultServiceClient interface {
	Test(ctx context.Context, in *TestReq, opts ...grpc.CallOption) (*TestRes, error)
}

type vaultServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVaultServiceClient(cc grpc.ClientConnInterface) VaultServiceClient {
	return &vaultServiceClient{cc}
}

func (c *vaultServiceClient) Test(ctx context.Context, in *TestReq, opts ...grpc.CallOption) (*TestRes, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TestRes)
	err := c.cc.Invoke(ctx, VaultService_Test_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VaultServiceServer is the server API for VaultService service.
// All implementations must embed UnimplementedVaultServiceServer
// for forward compatibility.
type VaultServiceServer interface {
	Test(context.Context, *TestReq) (*TestRes, error)
	mustEmbedUnimplementedVaultServiceServer()
}

// UnimplementedVaultServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedVaultServiceServer struct{}

func (UnimplementedVaultServiceServer) Test(context.Context, *TestReq) (*TestRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Test not implemented")
}
func (UnimplementedVaultServiceServer) mustEmbedUnimplementedVaultServiceServer() {}
func (UnimplementedVaultServiceServer) testEmbeddedByValue()                      {}

// UnsafeVaultServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VaultServiceServer will
// result in compilation errors.
type UnsafeVaultServiceServer interface {
	mustEmbedUnimplementedVaultServiceServer()
}

func RegisterVaultServiceServer(s grpc.ServiceRegistrar, srv VaultServiceServer) {
	// If the following call pancis, it indicates UnimplementedVaultServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&VaultService_ServiceDesc, srv)
}

func _VaultService_Test_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).Test(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VaultService_Test_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).Test(ctx, req.(*TestReq))
	}
	return interceptor(ctx, in, info, handler)
}

// VaultService_ServiceDesc is the grpc.ServiceDesc for VaultService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VaultService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "VaultService",
	HandlerType: (*VaultServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Test",
			Handler:    _VaultService_Test_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/vault.proto",
}
