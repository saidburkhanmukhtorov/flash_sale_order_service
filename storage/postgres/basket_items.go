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

type BasketItemRepo struct {
	db *pgx.Conn
}

func NewBasketItemRepo(db *pgx.Conn) *BasketItemRepo {
	return &BasketItemRepo{
		db: db,
	}
}

func (r *BasketItemRepo) CreateBasketItem(ctx context.Context, req *order_service.CreateBasketItemRequest) (*order_service.BasketItem, error) {
	if req.BasketItem.Id == "" {
		req.BasketItem.Id = uuid.NewString()
	}

	query := `
		INSERT INTO basket_items (
			id,
			basket_id,
			product_id,
			flash_sale_event_product_id,
			discount_product_id,
			quantity,
			unit_price,
			total_price,
			product_type,
			created_at,
			updated_at,
			deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW(), 0
		) RETURNING id, created_at, updated_at
	`

	flashSaleEventProductID := sql.NullString{
		String: req.BasketItem.FlashSaleEventProductId,
		Valid:  req.BasketItem.FlashSaleEventProductId != "",
	}
	discountProductID := sql.NullString{
		String: req.BasketItem.DiscountProductId,
		Valid:  req.BasketItem.DiscountProductId != "",
	}
	basketItemModel := makeBasketItemModel(req.BasketItem)
	err := r.db.QueryRow(ctx, query,
		basketItemModel.Id,
		basketItemModel.BasketId,
		basketItemModel.ProductId,
		flashSaleEventProductID, // Pass NullString here
		discountProductID,       // Pass NullString here
		basketItemModel.Quantity,
		basketItemModel.UnitPrice,
		basketItemModel.TotalPrice,
		basketItemModel.ProductType,
	).Scan(&basketItemModel.Id, &basketItemModel.CreatedAt, &basketItemModel.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return makeBasketItemProto(basketItemModel), nil
}
func (r *BasketItemRepo) GetBasketItem(ctx context.Context, req *order_service.GetBasketItemRequest) (*order_service.BasketItem, error) {
	var (
		basketItemModel models.BasketItem
		// Use sql.NullString for nullable UUID fields
		flashSaleEventProductID sql.NullString
		discountProductID       sql.NullString
	)

	query := `
		SELECT 
			id,
			basket_id,
			product_id,
			flash_sale_event_product_id,
			discount_product_id,
			quantity,
			unit_price,
			total_price,
			product_type,
			created_at,
			updated_at,
			deleted_at
		FROM basket_items
		WHERE id = $1 AND deleted_at = 0
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&basketItemModel.Id,
		&basketItemModel.BasketId,
		&basketItemModel.ProductId,
		&flashSaleEventProductID, // Scan into NullString
		&discountProductID,       // Scan into NullString
		&basketItemModel.Quantity,
		&basketItemModel.UnitPrice,
		&basketItemModel.TotalPrice,
		&basketItemModel.ProductType,
		&basketItemModel.CreatedAt,
		&basketItemModel.UpdatedAt,
		&basketItemModel.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	// Set the string fields in the model if Valid is true
	if flashSaleEventProductID.Valid {
		basketItemModel.FlashSaleEventProductId = flashSaleEventProductID.String
	}
	if discountProductID.Valid {
		basketItemModel.DiscountProductId = discountProductID.String
	}

	return makeBasketItemProto(basketItemModel), nil
}
func (r *BasketItemRepo) DeleteBasketItem(ctx context.Context, req *order_service.DeleteBasketItemRequest) (*order_service.DeleteBasketItemResponse, error) {
	query := `
		UPDATE basket_items
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at = 0
	`

	_, err := r.db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		return nil, err
	}

	return &order_service.DeleteBasketItemResponse{
		Message: "Basket item deleted successfully",
	}, nil
}
func (r *BasketItemRepo) ListBasketItems(ctx context.Context, req *order_service.ListBasketItemsRequest) (*order_service.ListBasketItemsResponse, error) {
	var args []interface{}
	count := 1
	query := `
		SELECT 
			id,
			basket_id,
			product_id,
			flash_sale_event_product_id,
			discount_product_id,
			quantity,
			unit_price,
			total_price,
			product_type,
			created_at,
			updated_at,
			deleted_at
		FROM 
			basket_items
		WHERE 1=1 AND deleted_at = 0
	`

	filter := ""

	if req.BasketId != "" {
		filter += fmt.Sprintf(" AND basket_id = $%d", count)
		args = append(args, req.BasketId)
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

	totalCountQuery := "SELECT count(*) FROM basket_items WHERE 1=1 AND deleted_at = 0" + filter
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

	var basketItemList []*order_service.BasketItem

	for rows.Next() {
		var (
			basketItemModel         models.BasketItem
			flashSaleEventProductID sql.NullString
			discountProductID       sql.NullString
		)
		err = rows.Scan(
			&basketItemModel.Id,
			&basketItemModel.BasketId,
			&basketItemModel.ProductId,
			&flashSaleEventProductID, // Scan into NullString
			&discountProductID,       // Scan into NullString
			&basketItemModel.Quantity,
			&basketItemModel.UnitPrice,
			&basketItemModel.TotalPrice,
			&basketItemModel.ProductType,
			&basketItemModel.CreatedAt,
			&basketItemModel.UpdatedAt,
			&basketItemModel.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set the string fields in the model if Valid is true
		if flashSaleEventProductID.Valid {
			basketItemModel.FlashSaleEventProductId = flashSaleEventProductID.String
		}
		if discountProductID.Valid {
			basketItemModel.DiscountProductId = discountProductID.String
		}

		basketItemList = append(basketItemList, makeBasketItemProto(basketItemModel))
	}

	return &order_service.ListBasketItemsResponse{
		BasketItems: basketItemList,
		Total:       int32(totalCount),
	}, nil
}

// Convert db model to proto model
func makeBasketItemProto(item models.BasketItem) *order_service.BasketItem {
	return &order_service.BasketItem{
		Id:                      item.Id,
		BasketId:                item.BasketId,
		ProductId:               item.ProductId,
		FlashSaleEventProductId: item.FlashSaleEventProductId,
		DiscountProductId:       item.DiscountProductId,
		Quantity:                item.Quantity,
		UnitPrice:               item.UnitPrice,
		TotalPrice:              item.TotalPrice,
		ProductType:             item.ProductType,
		CreatedAt:               timestamppb.New(item.CreatedAt),
		UpdatedAt:               timestamppb.New(item.UpdatedAt),
	}
}

// Convert proto model to db model
func makeBasketItemModel(item *order_service.BasketItem) models.BasketItem {
	return models.BasketItem{
		Id:                      item.Id,
		BasketId:                item.BasketId,
		ProductId:               item.ProductId,
		FlashSaleEventProductId: item.FlashSaleEventProductId,
		DiscountProductId:       item.DiscountProductId,
		Quantity:                item.Quantity,
		UnitPrice:               item.UnitPrice,
		TotalPrice:              item.TotalPrice,
		ProductType:             item.ProductType,
		CreatedAt:               item.CreatedAt.AsTime(),
		UpdatedAt:               item.UpdatedAt.AsTime(),
	}
}
