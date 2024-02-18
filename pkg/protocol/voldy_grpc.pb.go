// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: pkg/protocol/voldy.proto

package protocol

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

// VoldyClient is the client API for Voldy service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VoldyClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error)
}

type voldyClient struct {
	cc grpc.ClientConnInterface
}

func NewVoldyClient(cc grpc.ClientConnInterface) VoldyClient {
	return &voldyClient{cc}
}

func (c *voldyClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/Voldy/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *voldyClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error) {
	out := new(PutResponse)
	err := c.cc.Invoke(ctx, "/Voldy/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VoldyServer is the server API for Voldy service.
// All implementations must embed UnimplementedVoldyServer
// for forward compatibility
type VoldyServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Put(context.Context, *PutRequest) (*PutResponse, error)
	mustEmbedUnimplementedVoldyServer()
}

// UnimplementedVoldyServer must be embedded to have forward compatible implementations.
type UnimplementedVoldyServer struct {
}

func (UnimplementedVoldyServer) Get(context.Context, *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedVoldyServer) Put(context.Context, *PutRequest) (*PutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedVoldyServer) mustEmbedUnimplementedVoldyServer() {}

// UnsafeVoldyServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VoldyServer will
// result in compilation errors.
type UnsafeVoldyServer interface {
	mustEmbedUnimplementedVoldyServer()
}

func RegisterVoldyServer(s grpc.ServiceRegistrar, srv VoldyServer) {
	s.RegisterService(&Voldy_ServiceDesc, srv)
}

func _Voldy_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VoldyServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Voldy/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VoldyServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Voldy_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VoldyServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Voldy/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VoldyServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Voldy_ServiceDesc is the grpc.ServiceDesc for Voldy service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Voldy_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Voldy",
	HandlerType: (*VoldyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Voldy_Get_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _Voldy_Put_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/protocol/voldy.proto",
}