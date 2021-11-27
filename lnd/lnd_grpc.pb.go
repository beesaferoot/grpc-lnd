// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package lnd

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

// LNDClient is the client API for LND service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LNDClient interface {
	SpawnNodes(ctx context.Context, opts ...grpc.CallOption) (LND_SpawnNodesClient, error)
	GetNodesListByStatus(ctx context.Context, in *Status, opts ...grpc.CallOption) (*NodeList, error)
	DestroyNode(ctx context.Context, in *NodeId, opts ...grpc.CallOption) (*NodeDetail, error)
}

type lNDClient struct {
	cc grpc.ClientConnInterface
}

func NewLNDClient(cc grpc.ClientConnInterface) LNDClient {
	return &lNDClient{cc}
}

func (c *lNDClient) SpawnNodes(ctx context.Context, opts ...grpc.CallOption) (LND_SpawnNodesClient, error) {
	stream, err := c.cc.NewStream(ctx, &LND_ServiceDesc.Streams[0], "/LND/SpawnNodes", opts...)
	if err != nil {
		return nil, err
	}
	x := &lNDSpawnNodesClient{stream}
	return x, nil
}

type LND_SpawnNodesClient interface {
	Send(*NodeDetail) error
	CloseAndRecv() (*NodeList, error)
	grpc.ClientStream
}

type lNDSpawnNodesClient struct {
	grpc.ClientStream
}

func (x *lNDSpawnNodesClient) Send(m *NodeDetail) error {
	return x.ClientStream.SendMsg(m)
}

func (x *lNDSpawnNodesClient) CloseAndRecv() (*NodeList, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(NodeList)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *lNDClient) GetNodesListByStatus(ctx context.Context, in *Status, opts ...grpc.CallOption) (*NodeList, error) {
	out := new(NodeList)
	err := c.cc.Invoke(ctx, "/LND/GetNodesListByStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lNDClient) DestroyNode(ctx context.Context, in *NodeId, opts ...grpc.CallOption) (*NodeDetail, error) {
	out := new(NodeDetail)
	err := c.cc.Invoke(ctx, "/LND/DestroyNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LNDServer is the server API for LND service.
// All implementations must embed UnimplementedLNDServer
// for forward compatibility
type LNDServer interface {
	SpawnNodes(LND_SpawnNodesServer) error
	GetNodesListByStatus(context.Context, *Status) (*NodeList, error)
	DestroyNode(context.Context, *NodeId) (*NodeDetail, error)
	mustEmbedUnimplementedLNDServer()
}

// UnimplementedLNDServer must be embedded to have forward compatible implementations.
type UnimplementedLNDServer struct {
}

func (UnimplementedLNDServer) SpawnNodes(LND_SpawnNodesServer) error {
	return status.Errorf(codes.Unimplemented, "method SpawnNodes not implemented")
}
func (UnimplementedLNDServer) GetNodesListByStatus(context.Context, *Status) (*NodeList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodesListByStatus not implemented")
}
func (UnimplementedLNDServer) DestroyNode(context.Context, *NodeId) (*NodeDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DestroyNode not implemented")
}
func (UnimplementedLNDServer) mustEmbedUnimplementedLNDServer() {}

// UnsafeLNDServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LNDServer will
// result in compilation errors.
type UnsafeLNDServer interface {
	mustEmbedUnimplementedLNDServer()
}

func RegisterLNDServer(s grpc.ServiceRegistrar, srv LNDServer) {
	s.RegisterService(&LND_ServiceDesc, srv)
}

func _LND_SpawnNodes_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LNDServer).SpawnNodes(&lNDSpawnNodesServer{stream})
}

type LND_SpawnNodesServer interface {
	SendAndClose(*NodeList) error
	Recv() (*NodeDetail, error)
	grpc.ServerStream
}

type lNDSpawnNodesServer struct {
	grpc.ServerStream
}

func (x *lNDSpawnNodesServer) SendAndClose(m *NodeList) error {
	return x.ServerStream.SendMsg(m)
}

func (x *lNDSpawnNodesServer) Recv() (*NodeDetail, error) {
	m := new(NodeDetail)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _LND_GetNodesListByStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Status)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LNDServer).GetNodesListByStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/LND/GetNodesListByStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LNDServer).GetNodesListByStatus(ctx, req.(*Status))
	}
	return interceptor(ctx, in, info, handler)
}

func _LND_DestroyNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LNDServer).DestroyNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/LND/DestroyNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LNDServer).DestroyNode(ctx, req.(*NodeId))
	}
	return interceptor(ctx, in, info, handler)
}

// LND_ServiceDesc is the grpc.ServiceDesc for LND service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LND_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "LND",
	HandlerType: (*LNDServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNodesListByStatus",
			Handler:    _LND_GetNodesListByStatus_Handler,
		},
		{
			MethodName: "DestroyNode",
			Handler:    _LND_DestroyNode_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SpawnNodes",
			Handler:       _LND_SpawnNodes_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "lnd/lnd.proto",
}
