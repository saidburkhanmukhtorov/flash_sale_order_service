// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: submodule/order_service/order_items.proto

package order_service

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
	OrderItemService_GetOrderItem_FullMethodName              = "/order_service.OrderItemService/GetOrderItem"
	OrderItemService_ListOrderItems_FullMethodName            = "/order_service.OrderItemService/ListOrderItems"
	OrderItemService_ConvertBasketToOrderItems_FullMethodName = "/order_service.OrderItemService/ConvertBasketToOrderItems"
	OrderItemService_DeleteOrderItem_FullMethodName           = "/order_service.OrderItemService/DeleteOrderItem"
)

// OrderItemServiceClient is the client API for OrderItemService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// OrderItemService defines the gRPC service for managing order items.
type OrderItemServiceClient interface {
	GetOrderItem(ctx context.Context, in *GetOrderItemRequest, opts ...grpc.CallOption) (*GetOrderItemResponse, error)
	ListOrderItems(ctx context.Context, in *ListOrderItemsRequest, opts ...grpc.CallOption) (*ListOrderItemsResponse, error)
	ConvertBasketToOrderItems(ctx context.Context, in *ConvertBasketToOrderItemsRequest, opts ...grpc.CallOption) (*ConvertBasketToOrderItemsResponse, error)
	DeleteOrderItem(ctx context.Context, in *DeleteOrderItemRequest, opts ...grpc.CallOption) (*DeleteOrderItemResponse, error)
}

type orderItemServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderItemServiceClient(cc grpc.ClientConnInterface) OrderItemServiceClient {
	return &orderItemServiceClient{cc}
}

func (c *orderItemServiceClient) GetOrderItem(ctx context.Context, in *GetOrderItemRequest, opts ...grpc.CallOption) (*GetOrderItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetOrderItemResponse)
	err := c.cc.Invoke(ctx, OrderItemService_GetOrderItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderItemServiceClient) ListOrderItems(ctx context.Context, in *ListOrderItemsRequest, opts ...grpc.CallOption) (*ListOrderItemsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListOrderItemsResponse)
	err := c.cc.Invoke(ctx, OrderItemService_ListOrderItems_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderItemServiceClient) ConvertBasketToOrderItems(ctx context.Context, in *ConvertBasketToOrderItemsRequest, opts ...grpc.CallOption) (*ConvertBasketToOrderItemsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ConvertBasketToOrderItemsResponse)
	err := c.cc.Invoke(ctx, OrderItemService_ConvertBasketToOrderItems_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderItemServiceClient) DeleteOrderItem(ctx context.Context, in *DeleteOrderItemRequest, opts ...grpc.CallOption) (*DeleteOrderItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteOrderItemResponse)
	err := c.cc.Invoke(ctx, OrderItemService_DeleteOrderItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderItemServiceServer is the server API for OrderItemService service.
// All implementations must embed UnimplementedOrderItemServiceServer
// for forward compatibility.
//
// OrderItemService defines the gRPC service for managing order items.
type OrderItemServiceServer interface {
	GetOrderItem(context.Context, *GetOrderItemRequest) (*GetOrderItemResponse, error)
	ListOrderItems(context.Context, *ListOrderItemsRequest) (*ListOrderItemsResponse, error)
	ConvertBasketToOrderItems(context.Context, *ConvertBasketToOrderItemsRequest) (*ConvertBasketToOrderItemsResponse, error)
	DeleteOrderItem(context.Context, *DeleteOrderItemRequest) (*DeleteOrderItemResponse, error)
	mustEmbedUnimplementedOrderItemServiceServer()
}

// UnimplementedOrderItemServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOrderItemServiceServer struct{}

func (UnimplementedOrderItemServiceServer) GetOrderItem(context.Context, *GetOrderItemRequest) (*GetOrderItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrderItem not implemented")
}
func (UnimplementedOrderItemServiceServer) ListOrderItems(context.Context, *ListOrderItemsRequest) (*ListOrderItemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOrderItems not implemented")
}
func (UnimplementedOrderItemServiceServer) ConvertBasketToOrderItems(context.Context, *ConvertBasketToOrderItemsRequest) (*ConvertBasketToOrderItemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConvertBasketToOrderItems not implemented")
}
func (UnimplementedOrderItemServiceServer) DeleteOrderItem(context.Context, *DeleteOrderItemRequest) (*DeleteOrderItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteOrderItem not implemented")
}
func (UnimplementedOrderItemServiceServer) mustEmbedUnimplementedOrderItemServiceServer() {}
func (UnimplementedOrderItemServiceServer) testEmbeddedByValue()                          {}

// UnsafeOrderItemServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderItemServiceServer will
// result in compilation errors.
type UnsafeOrderItemServiceServer interface {
	mustEmbedUnimplementedOrderItemServiceServer()
}

func RegisterOrderItemServiceServer(s grpc.ServiceRegistrar, srv OrderItemServiceServer) {
	// If the following call pancis, it indicates UnimplementedOrderItemServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&OrderItemService_ServiceDesc, srv)
}

func _OrderItemService_GetOrderItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOrderItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderItemServiceServer).GetOrderItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderItemService_GetOrderItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderItemServiceServer).GetOrderItem(ctx, req.(*GetOrderItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderItemService_ListOrderItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListOrderItemsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderItemServiceServer).ListOrderItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderItemService_ListOrderItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderItemServiceServer).ListOrderItems(ctx, req.(*ListOrderItemsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderItemService_ConvertBasketToOrderItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConvertBasketToOrderItemsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderItemServiceServer).ConvertBasketToOrderItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderItemService_ConvertBasketToOrderItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderItemServiceServer).ConvertBasketToOrderItems(ctx, req.(*ConvertBasketToOrderItemsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderItemService_DeleteOrderItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteOrderItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderItemServiceServer).DeleteOrderItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OrderItemService_DeleteOrderItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderItemServiceServer).DeleteOrderItem(ctx, req.(*DeleteOrderItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderItemService_ServiceDesc is the grpc.ServiceDesc for OrderItemService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderItemService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order_service.OrderItemService",
	HandlerType: (*OrderItemServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOrderItem",
			Handler:    _OrderItemService_GetOrderItem_Handler,
		},
		{
			MethodName: "ListOrderItems",
			Handler:    _OrderItemService_ListOrderItems_Handler,
		},
		{
			MethodName: "ConvertBasketToOrderItems",
			Handler:    _OrderItemService_ConvertBasketToOrderItems_Handler,
		},
		{
			MethodName: "DeleteOrderItem",
			Handler:    _OrderItemService_DeleteOrderItem_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "submodule/order_service/order_items.proto",
}
