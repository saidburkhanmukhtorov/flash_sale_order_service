package service

import (
	"context"
	"fmt"
	"log"

	"github.com/flash_sale/flash_sale_order_service/genproto/order_service"
	"github.com/flash_sale/flash_sale_order_service/storage"
	"github.com/flash_sale/flash_sale_order_service/storage/redis"
)

// OrderItemService implements the order_service.OrderItemServiceServer interface.
type OrderItemService struct {
	storage     storage.StorageI
	redisClient *redis.Client
	order_service.UnimplementedOrderItemServiceServer
}

// NewOrderItemService creates a new OrderItemService instance.
func NewOrderItemService(storage storage.StorageI, redisClient *redis.Client) *OrderItemService {
	return &OrderItemService{
		storage:     storage,
		redisClient: redisClient,
	}
}

// GetOrderItem retrieves an order item by its ID.
func (s *OrderItemService) GetOrderItem(ctx context.Context, req *order_service.GetOrderItemRequest) (*order_service.GetOrderItemResponse, error) {
	orderItem, err := s.storage.OrderItem().GetOrderItem(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order item: %w", err)
	}

	return &order_service.GetOrderItemResponse{
		OrderItem: orderItem,
	}, nil
}

// ListOrderItems retrieves a list of order items.
func (s *OrderItemService) ListOrderItems(ctx context.Context, req *order_service.ListOrderItemsRequest) (*order_service.ListOrderItemsResponse, error) {
	response, err := s.storage.OrderItem().ListOrderItems(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list order items: %w", err)
	}

	return response, nil
}

// ConvertBasketToOrderItems converts basket items to order items and sends a notification.
func (s *OrderItemService) ConvertBasketToOrderItems(ctx context.Context, req *order_service.ConvertBasketToOrderItemsRequest) (*order_service.ConvertBasketToOrderItemsResponse, error) {
	orderID, err := s.storage.OrderItem().ConvertBasketToOrderItems(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to convert basket items to order items: %w", err)
	}

	// Get order details for notification
	order, err := s.storage.Order().GetOrder(ctx, &order_service.GetOrderRequest{Id: orderID.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to get order for notification: %w", err)
	}

	// Send notification to the user
	notificationMessage := fmt.Sprintf("Your order #%s is being prepared! We'll notify you when it's ready for pickup.", order.Id)
	if err := s.redisClient.AddNotification(ctx, order.ClientId, notificationMessage); err != nil {
		log.Printf("failed to send notification: %v", err)
		// Handle error (e.g., log and continue, retry, etc.)
	}

	return orderID, nil
}

// DeleteOrderItem deletes an order item by its ID and updates the order total price.
func (s *OrderItemService) DeleteOrderItem(ctx context.Context, req *order_service.DeleteOrderItemRequest) (*order_service.DeleteOrderItemResponse, error) {
	_, err := s.storage.OrderItem().DeleteOrderItem(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete order item: %w", err)
	}

	return &order_service.DeleteOrderItemResponse{
		Message: "Order item deleted successfully",
	}, nil
}
