package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/flash_sale/flash_sale_order_service/genproto/order_service"
	"github.com/flash_sale/flash_sale_order_service/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderItemRepo struct {
	db *pgx.Conn
}

func NewOrderItemRepo(db *pgx.Conn) *OrderItemRepo {
	return &OrderItemRepo{
		db: db,
	}
}
func (r *OrderItemRepo) GetOrderItem(ctx context.Context, req *order_service.GetOrderItemRequest) (*order_service.OrderItem, error) {
	var (
		orderItemModel          models.OrderItem
		flashSaleEventProductID sql.NullString
		discountProductID       sql.NullString
	)

	query := `
		SELECT 
			id,
			order_id,
			product_id,
			flash_sale_event_product_id,
			discount_product_id,
			quantity,
			unit_price,
			total_price,
			discount_applied,
			product_type,
			created_at,
			updated_at,
			deleted_at
		FROM order_items
		WHERE id = $1 AND deleted_at = 0
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&orderItemModel.Id,
		&orderItemModel.OrderId,
		&orderItemModel.ProductId,
		&flashSaleEventProductID,
		&discountProductID,
		&orderItemModel.Quantity,
		&orderItemModel.UnitPrice,
		&orderItemModel.TotalPrice,
		&orderItemModel.DiscountApplied,
		&orderItemModel.ProductType,
		&orderItemModel.CreatedAt,
		&orderItemModel.UpdatedAt,
		&orderItemModel.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	// Set the string fields in the model if Valid is true
	if flashSaleEventProductID.Valid {
		orderItemModel.FlashSaleEventProductId = flashSaleEventProductID.String
	}
	if discountProductID.Valid {
		orderItemModel.DiscountProductId = discountProductID.String
	}

	return makeOrderItemProto(orderItemModel), nil
}

func (r *OrderItemRepo) ListOrderItems(ctx context.Context, req *order_service.ListOrderItemsRequest) (*order_service.ListOrderItemsResponse, error) {
	var args []interface{}
	count := 1
	query := `
		SELECT 
			id,
			order_id,
			product_id,
			flash_sale_event_product_id,
			discount_product_id,
			quantity,
			unit_price,
			total_price,
			discount_applied,
			product_type,
			created_at,
			updated_at,
			deleted_at
		FROM 
			order_items
		WHERE 1=1 AND deleted_at = 0
	`

	filter := ""

	if req.OrderId != "" {
		filter += fmt.Sprintf(" AND order_id = $%d", count)
		args = append(args, req.OrderId)
		count++
	}

	query += filter

	// Handle invalid page or limit values
	if req.Page <= 0 {
		req.Page = 1 // Default to page 1
	}
	if req.Limit <= 0 {
		req.Limit = 10 // Default to a limit of 10
	}

	totalCountQuery := "SELECT count(*) FROM order_items WHERE 1=1 AND deleted_at = 0" + filter
	var totalCount int
	err := r.db.QueryRow(ctx, totalCountQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Add LIMIT and OFFSET for pagination using the proto fields
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", count, count+1)
	args = append(args, req.Limit, (req.Page-1)*req.Limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItemList []*order_service.OrderItem

	for rows.Next() {
		var (
			orderItemModel          models.OrderItem
			flashSaleEventProductID sql.NullString
			discountProductID       sql.NullString
		)
		err = rows.Scan(
			&orderItemModel.Id,
			&orderItemModel.OrderId,
			&orderItemModel.ProductId,
			&flashSaleEventProductID,
			&discountProductID,
			&orderItemModel.Quantity,
			&orderItemModel.UnitPrice,
			&orderItemModel.TotalPrice,
			&orderItemModel.DiscountApplied,
			&orderItemModel.ProductType,
			&orderItemModel.CreatedAt,
			&orderItemModel.UpdatedAt,
			&orderItemModel.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set the string fields in the model if Valid is true
		if flashSaleEventProductID.Valid {
			orderItemModel.FlashSaleEventProductId = flashSaleEventProductID.String
		}
		if discountProductID.Valid {
			orderItemModel.DiscountProductId = discountProductID.String
		}

		orderItemList = append(orderItemList, makeOrderItemProto(orderItemModel))
	}

	return &order_service.ListOrderItemsResponse{
		OrderItems: orderItemList,
		Total:      int32(totalCount),
	}, nil
}
func (r *OrderItemRepo) ConvertBasketToOrderItems(ctx context.Context, req *order_service.ConvertBasketToOrderItemsRequest) (*order_service.ConvertBasketToOrderItemsResponse, error) {
	rBasket := NewBasketItemRepo(r.db)
	// 1. Get basket items using ListBasketItems
	basketItemsResponse, err := rBasket.ListBasketItems(ctx, &order_service.ListBasketItemsRequest{
		BasketId: req.BasketId,
		Page:     1,   // Get all items in the basket
		Limit:    100, // Assuming a maximum of 100 items in a basket
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get basket items: %w", err)
	}
	basketItems := basketItemsResponse.BasketItems // Extract basket items from the response

	// 2. Create order items from basket items
	_, err = r.createOrderItemsFromBasketItems(ctx, req.OrderId, basketItems)
	if err != nil {
		return nil, fmt.Errorf("failed to create order items: %w", err)
	}

	// 3. Update order total price
	if err := r.updateOrderTotalPrice(ctx, req.OrderId); err != nil {
		return nil, fmt.Errorf("failed to update order total price: %w", err)
	}

	return &order_service.ConvertBasketToOrderItemsResponse{
		Id: req.OrderId,
	}, nil
}
func (r *OrderItemRepo) createOrderItemsFromBasketItems(ctx context.Context, orderID string, basketItems []*order_service.BasketItem) ([]*order_service.OrderItem, error) {

	// Maps to store validation results for discounts and flash sale events
	discountCache := make(map[string]bool)
	flashSaleEventCache := make(map[string]bool)

	var orderItems []*order_service.OrderItem

	for _, basketItem := range basketItems {
		// 1. Get product details
		var product models.Product
		err := r.db.QueryRow(ctx, `
            SELECT id, name, base_price, current_price, image_url, stock_quantity, created_at, updated_at
            FROM products
            WHERE id = $1 AND deleted_at = 0
        `, basketItem.ProductId).Scan(
			&product.Id,
			&product.Name,
			&product.BasePrice,
			&product.CurrentPrice,
			&product.ImageUrl,
			&product.StockQuantity,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get product: %w", err)
		}

		// 2. Calculate unit price based on product type and validity of discounts/flash sales
		unitPrice := product.BasePrice // Default to base price
		discountApplied := float32(0)

		switch basketItem.ProductType {
		case "REGULAR":
			// No discount or flash sale, use base price
		case "FLASH_SALE":
			if basketItem.FlashSaleEventProductId != "" {
				isValid, ok := flashSaleEventCache[basketItem.FlashSaleEventProductId]
				if !ok {
					// Check if flash sale event is valid and cache the result
					isValid = isFlashSaleEventProductValid(ctx, r.db, basketItem.FlashSaleEventProductId)
					flashSaleEventCache[basketItem.FlashSaleEventProductId] = isValid
				}

				if isValid {
					var flashSaleEventProduct models.FlashSaleEventProduct
					err = r.db.QueryRow(ctx, `
                        SELECT sale_price
                        FROM flash_sale_event_products
                        WHERE id = $1 AND deleted_at = 0
                    `, basketItem.FlashSaleEventProductId).Scan(
						&flashSaleEventProduct.SalePrice,
					)
					if err != nil {
						return nil, fmt.Errorf("failed to get flash sale event product: %w", err)
					}

					unitPrice = flashSaleEventProduct.SalePrice
					discountApplied = product.BasePrice - unitPrice
				} // else use base price
			}
		case "DISCOUNT":
			if basketItem.DiscountProductId != "" {
				isValid, ok := discountCache[basketItem.DiscountProductId]
				if !ok {
					// Check if discount is valid and cache the result
					isValid = isDiscountValid(ctx, r.db, basketItem.DiscountProductId)
					discountCache[basketItem.DiscountProductId] = isValid
				}

				if isValid {
					var discount models.Discount
					err = r.db.QueryRow(ctx, `
                        SELECT discount_type, discount_value
                        FROM discounts
                        WHERE id = $1 AND deleted_at = 0
                    `, basketItem.DiscountProductId).Scan(
						&discount.DiscountType,
						&discount.DiscountValue,
					)
					if err != nil {
						return nil, fmt.Errorf("failed to get discount: %w", err)
					}

					unitPrice = calculateDiscountedPrice(product.BasePrice, &discount)
					discountApplied = product.BasePrice - unitPrice
				} // else use base price
			}
		}

		// 3. Create order item
		orderItem := &order_service.OrderItem{
			Id:                      uuid.NewString(),
			OrderId:                 orderID,
			ProductId:               basketItem.ProductId,
			FlashSaleEventProductId: basketItem.FlashSaleEventProductId,
			DiscountProductId:       basketItem.DiscountProductId,
			Quantity:                basketItem.Quantity,
			UnitPrice:               unitPrice,
			TotalPrice:              unitPrice * float32(basketItem.Quantity),
			DiscountApplied:         discountApplied,
			ProductType:             basketItem.ProductType,
			CreatedAt:               timestamppb.Now(),
			UpdatedAt:               timestamppb.Now(),
		}

		// Use sql.NullString for nullable UUIDs
		flashSaleEventProductID := sql.NullString{
			String: orderItem.FlashSaleEventProductId,
			Valid:  orderItem.FlashSaleEventProductId != "",
		}
		discountProductID := sql.NullString{
			String: orderItem.DiscountProductId,
			Valid:  orderItem.DiscountProductId != "",
		}

		query := `
			INSERT INTO order_items (
				id,
				order_id,
				product_id,
				flash_sale_event_product_id,
				discount_product_id,
				quantity,
				unit_price,
				total_price,
				discount_applied,
				product_type,
				created_at,
				updated_at,
				deleted_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW(), 0
			)
		`

		_, err = r.db.Exec(ctx, query,
			orderItem.Id,
			orderItem.OrderId,
			orderItem.ProductId,
			flashSaleEventProductID, // Pass NullString here
			discountProductID,       // Pass NullString here
			orderItem.Quantity,
			orderItem.UnitPrice,
			orderItem.TotalPrice,
			orderItem.DiscountApplied,
			orderItem.ProductType,
		)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}
func (r *OrderItemRepo) DeleteOrderItem(ctx context.Context, req *order_service.DeleteOrderItemRequest) (string, error) {
	// Get the order ID before deleting the order item
	var orderID string
	err := r.db.QueryRow(ctx, `
		SELECT order_id
		FROM order_items
		WHERE id = $1
	`, req.Id).Scan(&orderID)
	if err != nil {
		return "", fmt.Errorf("failed to get order ID: %w", err)
	}

	query := `
		DELETE FROM order_items
		WHERE id = $1
	`

	_, err = r.db.Exec(ctx, query, req.Id)
	if err != nil {
		return "", err
	}

	// Update the order's total price after deleting the item
	if err := r.updateOrderTotalPrice(ctx, orderID); err != nil {
		return "", fmt.Errorf("failed to update order total price: %w", err)
	}

	return orderID, nil
}

// Helper function to check if a flash sale event product is valid
func isFlashSaleEventProductValid(ctx context.Context, db *pgx.Conn, flashSaleEventProductID string) bool {
	var (
		status  string
		endTime time.Time
	)
	query := `
		SELECT fse.status, fse.end_time
		FROM flash_sale_event_products fsep
		JOIN flash_sale_events fse ON fsep.event_id = fse.id
		WHERE fsep.id = $1 AND fsep.deleted_at = 0
	`
	err := db.QueryRow(ctx, query, flashSaleEventProductID).Scan(&status, &endTime)
	if err != nil {
		return false // Handle the error appropriately
	}

	return status == "ACTIVE" && endTime.After(time.Now())
}

// Helper function to check if a discount is valid
func isDiscountValid(ctx context.Context, db *pgx.Conn, discountID string) bool {
	var (
		isActive bool
		endDate  time.Time
	)
	query := `
		SELECT is_active, end_date
		FROM discounts
		WHERE id = $1 AND deleted_at = 0
	`
	err := db.QueryRow(ctx, query, discountID).Scan(&isActive, &endDate)
	if err != nil {
		return false // Handle the error appropriately
	}

	return isActive && endDate.After(time.Now())
}

// Convert db model to proto model
func makeOrderItemProto(item models.OrderItem) *order_service.OrderItem {
	return &order_service.OrderItem{
		Id:                      item.Id,
		OrderId:                 item.OrderId,
		ProductId:               item.ProductId,
		FlashSaleEventProductId: item.FlashSaleEventProductId,
		DiscountProductId:       item.DiscountProductId,
		Quantity:                item.Quantity,
		UnitPrice:               item.UnitPrice,
		TotalPrice:              item.TotalPrice,
		DiscountApplied:         item.DiscountApplied,
		ProductType:             item.ProductType,
		CreatedAt:               timestamppb.New(item.CreatedAt),
		UpdatedAt:               timestamppb.New(item.UpdatedAt),
	}
}

// Helper function to calculate the discounted price
func calculateDiscountedPrice(basePrice float32, discount *models.Discount) float32 {
	if discount.DiscountType == "PERCENTAGE" {
		return basePrice * (1 - discount.DiscountValue/100)
	} else if discount.DiscountType == "FIXED_AMOUNT" {
		return basePrice - discount.DiscountValue
	}
	return basePrice
}

// Helper function to update the order's total price
func (r *OrderItemRepo) updateOrderTotalPrice(ctx context.Context, orderID string) error {
	var totalPrice sql.NullFloat64
	err := r.db.QueryRow(ctx, `
		SELECT SUM(total_price)
		FROM order_items
		WHERE order_id = $1 AND deleted_at = 0
	`, orderID).Scan(&totalPrice)
	if err != nil {
		return err
	}

	query := `
		UPDATE orders
		SET total_price = $1
		WHERE id = $2
	`

	_, err = r.db.Exec(ctx, query, totalPrice.Float64, orderID)
	return err
}
