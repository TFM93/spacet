// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: spacet.proto

package v1

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
	SpaceTService_LaunchBooking_FullMethodName = "/spacet.v1.SpaceTService/LaunchBooking"
	SpaceTService_CancelBooking_FullMethodName = "/spacet.v1.SpaceTService/CancelBooking"
	SpaceTService_ListBookings_FullMethodName  = "/spacet.v1.SpaceTService/ListBookings"
)

// SpaceTServiceClient is the client API for SpaceTService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// SpaceT service definition
type SpaceTServiceClient interface {
	LaunchBooking(ctx context.Context, in *BookingRequest, opts ...grpc.CallOption) (*Ticket, error)
	CancelBooking(ctx context.Context, in *TicketID, opts ...grpc.CallOption) (*TicketID, error)
	ListBookings(ctx context.Context, in *ListTicketsRequest, opts ...grpc.CallOption) (*ListTicketsResponse, error)
}

type spaceTServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSpaceTServiceClient(cc grpc.ClientConnInterface) SpaceTServiceClient {
	return &spaceTServiceClient{cc}
}

func (c *spaceTServiceClient) LaunchBooking(ctx context.Context, in *BookingRequest, opts ...grpc.CallOption) (*Ticket, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Ticket)
	err := c.cc.Invoke(ctx, SpaceTService_LaunchBooking_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceTServiceClient) CancelBooking(ctx context.Context, in *TicketID, opts ...grpc.CallOption) (*TicketID, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TicketID)
	err := c.cc.Invoke(ctx, SpaceTService_CancelBooking_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceTServiceClient) ListBookings(ctx context.Context, in *ListTicketsRequest, opts ...grpc.CallOption) (*ListTicketsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListTicketsResponse)
	err := c.cc.Invoke(ctx, SpaceTService_ListBookings_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SpaceTServiceServer is the server API for SpaceTService service.
// All implementations must embed UnimplementedSpaceTServiceServer
// for forward compatibility.
//
// SpaceT service definition
type SpaceTServiceServer interface {
	LaunchBooking(context.Context, *BookingRequest) (*Ticket, error)
	CancelBooking(context.Context, *TicketID) (*TicketID, error)
	ListBookings(context.Context, *ListTicketsRequest) (*ListTicketsResponse, error)
	mustEmbedUnimplementedSpaceTServiceServer()
}

// UnimplementedSpaceTServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSpaceTServiceServer struct{}

func (UnimplementedSpaceTServiceServer) LaunchBooking(context.Context, *BookingRequest) (*Ticket, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LaunchBooking not implemented")
}
func (UnimplementedSpaceTServiceServer) CancelBooking(context.Context, *TicketID) (*TicketID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelBooking not implemented")
}
func (UnimplementedSpaceTServiceServer) ListBookings(context.Context, *ListTicketsRequest) (*ListTicketsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBookings not implemented")
}
func (UnimplementedSpaceTServiceServer) mustEmbedUnimplementedSpaceTServiceServer() {}
func (UnimplementedSpaceTServiceServer) testEmbeddedByValue()                       {}

// UnsafeSpaceTServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SpaceTServiceServer will
// result in compilation errors.
type UnsafeSpaceTServiceServer interface {
	mustEmbedUnimplementedSpaceTServiceServer()
}

func RegisterSpaceTServiceServer(s grpc.ServiceRegistrar, srv SpaceTServiceServer) {
	// If the following call pancis, it indicates UnimplementedSpaceTServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SpaceTService_ServiceDesc, srv)
}

func _SpaceTService_LaunchBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceTServiceServer).LaunchBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpaceTService_LaunchBooking_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceTServiceServer).LaunchBooking(ctx, req.(*BookingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceTService_CancelBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TicketID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceTServiceServer).CancelBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpaceTService_CancelBooking_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceTServiceServer).CancelBooking(ctx, req.(*TicketID))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceTService_ListBookings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTicketsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceTServiceServer).ListBookings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpaceTService_ListBookings_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceTServiceServer).ListBookings(ctx, req.(*ListTicketsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SpaceTService_ServiceDesc is the grpc.ServiceDesc for SpaceTService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SpaceTService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spacet.v1.SpaceTService",
	HandlerType: (*SpaceTServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LaunchBooking",
			Handler:    _SpaceTService_LaunchBooking_Handler,
		},
		{
			MethodName: "CancelBooking",
			Handler:    _SpaceTService_CancelBooking_Handler,
		},
		{
			MethodName: "ListBookings",
			Handler:    _SpaceTService_ListBookings_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "spacet.proto",
}
