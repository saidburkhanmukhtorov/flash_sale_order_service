package storage

import (
	"context"

	"github.com/flash_sale/flash_sale_order_service/genproto/order_service"
)

// StorageI defines the interface for interacting with the storage layer.
type StorageI interface {
	Basket() BasketI
	BasketItem() BasketItemI
	Order() OrderI
	OrderItem() OrderItemI
}

// BasketI defines methods for interacting with basket data.
type BasketI interface {
	CreateBasket(ctx context.Context, req *order_service.CreateBasketRequest) (*order_service.Basket, error)
	GetBasket(ctx context.Context, req *order_service.GetBasketRequest) (*order_service.Basket, error)
	UpdateBasket(ctx context.Context, req *order_service.UpdateBasketRequest) (*order_service.Basket, error)
	DeleteBasket(ctx context.Context, req *order_service.DeleteBasketRequest) (*order_service.DeleteBasketResponse, error)
	ListBaskets(ctx context.Context, req *order_service.ListBasketsRequest) (*order_service.ListBasketsResponse, error)
	UpdateBasketStatus(ctx context.Context, req *order_service.UpdateBasketStatusRequest) (*order_service.Basket, error)
}

// BasketItemI defines methods for interacting with basket item data.
type BasketItemI interface {
	CreateBasketItem(ctx context.Context, req *order_service.CreateBasketItemRequest) (*order_service.BasketItem, error)
	GetBasketItem(ctx context.Context, req *order_service.GetBasketItemRequest) (*order_service.BasketItem, error)
	DeleteBasketItem(ctx context.Context, req *order_service.DeleteBasketItemRequest) (*order_service.DeleteBasketItemResponse, error)
	ListBasketItems(ctx context.Context, req *order_service.ListBasketItemsRequest) (*order_service.ListBasketItemsResponse, error)
}

// OrderI defines methods for interacting with order data.
type OrderI interface {
	CreateOrder(ctx context.Context, req *order_service.CreateOrderRequest) (*order_service.Order, error)
	GetOrder(ctx context.Context, req *order_service.GetOrderRequest) (*order_service.Order, error)
	UpdateOrder(ctx context.Context, req *order_service.UpdateOrderRequest) (*order_service.Order, error)
	DeleteOrder(ctx context.Context, req *order_service.DeleteOrderRequest) (*order_service.DeleteOrderResponse, error)
	ListOrders(ctx context.Context, req *order_service.ListOrdersRequest) (*order_service.ListOrdersResponse, error)
	UpdateOrderStatus(ctx context.Context, req *order_service.UpdateOrderStatusRequest) (*order_service.Order, error)
}

// OrderItemI defines methods for interacting with order item data.
type OrderItemI interface {
	GetOrderItem(ctx context.Context, req *order_service.GetOrderItemRequest) (*order_service.OrderItem, error)
	ListOrderItems(ctx context.Context, req *order_service.ListOrderItemsRequest) (*order_service.ListOrderItemsResponse, error)
	ConvertBasketToOrderItems(ctx context.Context, req *order_service.ConvertBasketToOrderItemsRequest) (*order_service.ConvertBasketToOrderItemsResponse, error)
	DeleteOrderItem(ctx context.Context, req *order_service.DeleteOrderItemRequest) (string, error)
}
