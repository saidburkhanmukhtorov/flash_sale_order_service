package service

import (
	"context"
	"fmt"
	"log"

	"github.com/flash_sale/flash_sale_order_service/genproto/order_service"
	"github.com/flash_sale/flash_sale_order_service/storage"
	"github.com/flash_sale/flash_sale_order_service/storage/redis"
)

// OrderService implements the order_service.OrderServiceServer interface.
type OrderService struct {
	storage     storage.StorageI
	redisClient *redis.Client
	order_service.UnimplementedOrderServiceServer
}

// NewOrderService creates a new OrderService instance.
func NewOrderService(storage storage.StorageI, redisClient *redis.Client) *OrderService {
	return &OrderService{
		storage:     storage,
		redisClient: redisClient,
	}
}

// CreateOrder creates a new order.
func (s *OrderService) CreateOrder(ctx context.Context, req *order_service.CreateOrderRequest) (*order_service.CreateOrderResponse, error) {
	order, err := s.storage.Order().CreateOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return &order_service.CreateOrderResponse{
		Order: order,
	}, nil
}

// GetOrder retrieves an order by its ID.
func (s *OrderService) GetOrder(ctx context.Context, req *order_service.GetOrderRequest) (*order_service.GetOrderResponse, error) {
	order, err := s.storage.Order().GetOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order_service.GetOrderResponse{
		Order: order,
	}, nil
}

// UpdateOrder updates an existing order.
func (s *OrderService) UpdateOrder(ctx context.Context, req *order_service.UpdateOrderRequest) (*order_service.UpdateOrderResponse, error) {
	order, err := s.storage.Order().UpdateOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return &order_service.UpdateOrderResponse{
		Order: order,
	}, nil
}

// DeleteOrder deletes an order by its ID.
func (s *OrderService) DeleteOrder(ctx context.Context, req *order_service.DeleteOrderRequest) (*order_service.DeleteOrderResponse, error) {
	response, err := s.storage.Order().DeleteOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete order: %w", err)
	}

	return response, nil
}

// ListOrders retrieves a list of orders.
func (s *OrderService) ListOrders(ctx context.Context, req *order_service.ListOrdersRequest) (*order_service.ListOrdersResponse, error) {
	response, err := s.storage.Order().ListOrders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	return response, nil
}

// UpdateOrderStatus updates the status of an order and sends a notification.
func (s *OrderService) UpdateOrderStatus(ctx context.Context, req *order_service.UpdateOrderStatusRequest) (*order_service.UpdateOrderStatusResponse, error) {
	order, err := s.storage.Order().UpdateOrderStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// Send notification to the user
	notificationMessage := fmt.Sprintf("Your order status has been updated to %s.", order.Status)
	if err := s.redisClient.AddNotification(ctx, order.ClientId, notificationMessage); err != nil {
		log.Printf("failed to send notification: %v", err)
		// Handle error (e.g., log and continue, retry, etc.)
	}

	return &order_service.UpdateOrderStatusResponse{
		Order: order,
	}, nil
}
