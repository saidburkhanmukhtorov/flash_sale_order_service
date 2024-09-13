package service

import (
	"context"
	"fmt"

	"github.com/flash_sale/flash_sale_order_service/genproto/order_service"
	"github.com/flash_sale/flash_sale_order_service/storage"
)

// BasketItemService implements the order_service.BasketItemServiceServer interface.
type BasketItemService struct {
	storage storage.StorageI
	order_service.UnimplementedBasketItemServiceServer
}

// NewBasketItemService creates a new BasketItemService instance.
func NewBasketItemService(storage storage.StorageI) *BasketItemService {
	return &BasketItemService{
		storage: storage,
	}
}

// CreateBasketItem creates a new basket item.
func (s *BasketItemService) CreateBasketItem(ctx context.Context, req *order_service.CreateBasketItemRequest) (*order_service.CreateBasketItemResponse, error) {
	basketItem, err := s.storage.BasketItem().CreateBasketItem(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create basket item: %w", err)
	}

	return &order_service.CreateBasketItemResponse{
		BasketItem: basketItem,
	}, nil
}

// GetBasketItem retrieves a basket item by its ID.
func (s *BasketItemService) GetBasketItem(ctx context.Context, req *order_service.GetBasketItemRequest) (*order_service.GetBasketItemResponse, error) {
	basketItem, err := s.storage.BasketItem().GetBasketItem(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get basket item: %w", err)
	}

	return &order_service.GetBasketItemResponse{
		BasketItem: basketItem,
	}, nil
}

// DeleteBasketItem deletes a basket item by its ID.
func (s *BasketItemService) DeleteBasketItem(ctx context.Context, req *order_service.DeleteBasketItemRequest) (*order_service.DeleteBasketItemResponse, error) {
	response, err := s.storage.BasketItem().DeleteBasketItem(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to delete basket item: %w", err)
	}

	return response, nil
}

// ListBasketItems retrieves a list of basket items.
func (s *BasketItemService) ListBasketItems(ctx context.Context, req *order_service.ListBasketItemsRequest) (*order_service.ListBasketItemsResponse, error) {
	response, err := s.storage.BasketItem().ListBasketItems(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list basket items: %w", err)
	}

	return response, nil
}
