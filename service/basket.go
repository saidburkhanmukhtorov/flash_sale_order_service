package service

import (
	"context"
	"fmt"
	"log"

	"github.com/flash_sale/flash_sale_order_service/genproto/order_service"
	"github.com/flash_sale/flash_sale_order_service/storage"
	"github.com/flash_sale/flash_sale_order_service/storage/redis"
)

// BasketService implements the order_service.BasketServiceServer interface.
type BasketService struct {
	storage     storage.StorageI
	redisClient *redis.Client
	order_service.UnimplementedBasketServiceServer
}

// NewBasketService creates a new BasketService instance.
func NewBasketService(storage storage.StorageI, redisClient *redis.Client) *BasketService {
	return &BasketService{
		storage:     storage,
		redisClient: redisClient,
	}
}

// CreateBasket creates a new basket.
func (s *BasketService) CreateBasket(ctx context.Context, req *order_service.CreateBasketRequest) (*order_service.CreateBasketResponse, error) {
	basket, err := s.storage.Basket().CreateBasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create basket: %w", err)
	}

	return &order_service.CreateBasketResponse{
		Basket: basket,
	}, nil
}

// GetBasket retrieves a basket by its ID.
func (s *BasketService) GetBasket(ctx context.Context, req *order_service.GetBasketRequest) (*order_service.GetBasketResponse, error) {
	basket, err := s.storage.Basket().GetBasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get basket: %w", err)
	}

	return &order_service.GetBasketResponse{
		Basket: basket,
	}, nil
}

// UpdateBasket updates an existing basket.
func (s *BasketService) UpdateBasket(ctx context.Context, req *order_service.UpdateBasketRequest) (*order_service.UpdateBasketResponse, error) {
	basket, err := s.storage.Basket().UpdateBasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update basket: %w", err)
	}

	return &order_service.UpdateBasketResponse{
		Basket: basket,
	}, nil
}

// DeleteBasket deletes a basket by its ID.
func (s *BasketService) DeleteBasket(ctx context.Context, req *order_service.DeleteBasketRequest) (*order_service.DeleteBasketResponse, error) {
	response, err := s.storage.Basket().DeleteBasket(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete basket: %w", err)
	}

	return response, nil
}

// ListBaskets retrieves a list of baskets.
func (s *BasketService) ListBaskets(ctx context.Context, req *order_service.ListBasketsRequest) (*order_service.ListBasketsResponse, error) {
	response, err := s.storage.Basket().ListBaskets(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list baskets: %w", err)
	}

	return response, nil
}

// UpdateBasketStatus updates the status of a basket and sends a notification.
func (s *BasketService) UpdateBasketStatus(ctx context.Context, req *order_service.UpdateBasketStatusRequest) (*order_service.UpdateBasketStatusResponse, error) {
	basket, err := s.storage.Basket().UpdateBasketStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update basket status: %w", err)
	}

	// Send notification to the user
	notificationMessage := fmt.Sprintf("Your basket status has been updated to %s.", basket.Status)
	if err := s.redisClient.AddNotification(ctx, basket.UserId, notificationMessage); err != nil {
		log.Printf("failed to send notification: %v", err)
		// Handle error (e.g., log and continue, retry, etc.)
	}

	return &order_service.UpdateBasketStatusResponse{
		Basket: basket,
	}, nil
}
