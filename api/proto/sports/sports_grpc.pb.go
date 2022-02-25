// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.3
// source: sports/sports.proto

package sports

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

// SportsClient is the client API for Sports service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SportsClient interface {
	// ListSports returns a list of all sports.
	ListSports(ctx context.Context, in *ListSportsRequest, opts ...grpc.CallOption) (*ListSportsResponse, error)
	// GetSportDetails returns details of a single sports event
	GetSportsEventDetails(ctx context.Context, in *GetSportsDetailsRequest, opts ...grpc.CallOption) (*GetSportsDetailsResponse, error)
}

type sportsClient struct {
	cc grpc.ClientConnInterface
}

func NewSportsClient(cc grpc.ClientConnInterface) SportsClient {
	return &sportsClient{cc}
}

func (c *sportsClient) ListSports(ctx context.Context, in *ListSportsRequest, opts ...grpc.CallOption) (*ListSportsResponse, error) {
	out := new(ListSportsResponse)
	err := c.cc.Invoke(ctx, "/sports.Sports/ListSports", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sportsClient) GetSportsEventDetails(ctx context.Context, in *GetSportsDetailsRequest, opts ...grpc.CallOption) (*GetSportsDetailsResponse, error) {
	out := new(GetSportsDetailsResponse)
	err := c.cc.Invoke(ctx, "/sports.Sports/GetSportsEventDetails", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SportsServer is the server API for Sports service.
// All implementations must embed UnimplementedSportsServer
// for forward compatibility
type SportsServer interface {
	// ListSports returns a list of all sports.
	ListSports(context.Context, *ListSportsRequest) (*ListSportsResponse, error)
	// GetSportDetails returns details of a single sports event
	GetSportsEventDetails(context.Context, *GetSportsDetailsRequest) (*GetSportsDetailsResponse, error)
	mustEmbedUnimplementedSportsServer()
}

// UnimplementedSportsServer must be embedded to have forward compatible implementations.
type UnimplementedSportsServer struct {
}

func (UnimplementedSportsServer) ListSports(context.Context, *ListSportsRequest) (*ListSportsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSports not implemented")
}
func (UnimplementedSportsServer) GetSportsEventDetails(context.Context, *GetSportsDetailsRequest) (*GetSportsDetailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSportsEventDetails not implemented")
}
func (UnimplementedSportsServer) mustEmbedUnimplementedSportsServer() {}

// UnsafeSportsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SportsServer will
// result in compilation errors.
type UnsafeSportsServer interface {
	mustEmbedUnimplementedSportsServer()
}

func RegisterSportsServer(s grpc.ServiceRegistrar, srv SportsServer) {
	s.RegisterService(&Sports_ServiceDesc, srv)
}

func _Sports_ListSports_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSportsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SportsServer).ListSports(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sports.Sports/ListSports",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SportsServer).ListSports(ctx, req.(*ListSportsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Sports_GetSportsEventDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSportsDetailsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SportsServer).GetSportsEventDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sports.Sports/GetSportsEventDetails",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SportsServer).GetSportsEventDetails(ctx, req.(*GetSportsDetailsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Sports_ServiceDesc is the grpc.ServiceDesc for Sports service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sports_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sports.Sports",
	HandlerType: (*SportsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListSports",
			Handler:    _Sports_ListSports_Handler,
		},
		{
			MethodName: "GetSportsEventDetails",
			Handler:    _Sports_GetSportsEventDetails_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sports/sports.proto",
}
