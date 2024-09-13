package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/flash_sale/flash_sale_order_service/genproto/order_service"
	"github.com/flash_sale/flash_sale_order_service/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderRepo struct {
	db *pgx.Conn
}

func NewOrderRepo(db *pgx.Conn) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (r *OrderRepo) CreateOrder(ctx context.Context, req *order_service.CreateOrderRequest) (*order_service.Order, error) {
	if req.Order.Id == "" {
		req.Order.Id = uuid.NewString()
	}

	query := `
		INSERT INTO orders (
			id,
			client_id,
			delivery_latitude,
			delivery_longitude,
			total_price,
			status,
			created_at,
			updated_at,
			deleted_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, NOW(), NOW(), 0
		) RETURNING id, created_at, updated_at
	`

	orderModel := makeOrderModel(req.Order)

	err := r.db.QueryRow(ctx, query,
		orderModel.Id,
		orderModel.ClientId,
		orderModel.DeliveryLatitude,
		orderModel.DeliveryLongitude,
		orderModel.TotalPrice,
		orderModel.Status,
	).Scan(&orderModel.Id, &orderModel.CreatedAt, &orderModel.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return makeOrderProto(orderModel), nil
}

func (r *OrderRepo) GetOrder(ctx context.Context, req *order_service.GetOrderRequest) (*order_service.Order, error) {
	var orderModel models.Order

	query := `
		SELECT 
			id,
			client_id,
			delivery_latitude,
			delivery_longitude,
			total_price,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM orders
		WHERE id = $1 AND deleted_at = 0
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&orderModel.Id,
		&orderModel.ClientId,
		&orderModel.DeliveryLatitude,
		&orderModel.DeliveryLongitude,
		&orderModel.TotalPrice,
		&orderModel.Status,
		&orderModel.CreatedAt,
		&orderModel.UpdatedAt,
		&orderModel.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeOrderProto(orderModel), nil
}

func (r *OrderRepo) UpdateOrder(ctx context.Context, req *order_service.UpdateOrderRequest) (*order_service.Order, error) {
	query := `
		UPDATE orders
		SET 
			client_id = $1,
			delivery_latitude = $2,
			delivery_longitude = $3,
			total_price = $4,
			status = $5,
			updated_at = NOW()
		WHERE id = $6 AND deleted_at = 0
		RETURNING id, client_id, delivery_latitude, delivery_longitude, total_price, status, created_at, updated_at
	`

	orderModel := makeOrderModel(req.Order)

	err := r.db.QueryRow(ctx, query,
		orderModel.ClientId,
		orderModel.DeliveryLatitude,
		orderModel.DeliveryLongitude,
		orderModel.TotalPrice,
		orderModel.Status,
		orderModel.Id,
	).Scan(
		&orderModel.Id,
		&orderModel.ClientId,
		&orderModel.DeliveryLatitude,
		&orderModel.DeliveryLongitude,
		&orderModel.TotalPrice,
		&orderModel.Status,
		&orderModel.CreatedAt,
		&orderModel.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeOrderProto(orderModel), nil
}

func (r *OrderRepo) DeleteOrder(ctx context.Context, req *order_service.DeleteOrderRequest) (*order_service.DeleteOrderResponse, error) {
	query := `
		UPDATE orders
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at = 0
	`

	_, err := r.db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		return nil, err
	}

	return &order_service.DeleteOrderResponse{
		Message: "Order deleted successfully",
	}, nil
}

func (r *OrderRepo) ListOrders(ctx context.Context, req *order_service.ListOrdersRequest) (*order_service.ListOrdersResponse, error) {
	var args []interface{}
	count := 1
	query := `
		SELECT 
			id,
			client_id,
			delivery_latitude,
			delivery_longitude,
			total_price,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM 
			orders
		WHERE 1=1 AND deleted_at = 0
	`

	filter := ""

	if req.ClientId != "" {
		filter += fmt.Sprintf(" AND client_id = $%d", count)
		args = append(args, req.ClientId)
		count++
	}

	if req.Status != "" {
		filter += fmt.Sprintf(" AND status = $%d", count)
		args = append(args, req.Status)
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

	totalCountQuery := "SELECT count(*) FROM orders WHERE 1=1 AND deleted_at = 0" + filter
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

	var orderList []*order_service.Order

	for rows.Next() {
		var orderModel models.Order
		err = rows.Scan(
			&orderModel.Id,
			&orderModel.ClientId,
			&orderModel.DeliveryLatitude,
			&orderModel.DeliveryLongitude,
			&orderModel.TotalPrice,
			&orderModel.Status,
			&orderModel.CreatedAt,
			&orderModel.UpdatedAt,
			&orderModel.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		orderList = append(orderList, makeOrderProto(orderModel))
	}

	return &order_service.ListOrdersResponse{
		Orders: orderList,
		Total:  int32(totalCount),
	}, nil
}

func (r *OrderRepo) UpdateOrderStatus(ctx context.Context, req *order_service.UpdateOrderStatusRequest) (*order_service.Order, error) {
	query := `
		UPDATE orders
		SET 
			status = $1,
			updated_at = NOW()
		WHERE id = $2 AND deleted_at = 0
		RETURNING id, client_id, delivery_latitude, delivery_longitude, total_price, status, created_at, updated_at
	`

	var orderModel models.Order

	err := r.db.QueryRow(ctx, query,
		req.Status,
		req.Id,
	).Scan(
		&orderModel.Id,
		&orderModel.ClientId,
		&orderModel.DeliveryLatitude,
		&orderModel.DeliveryLongitude,
		&orderModel.TotalPrice,
		&orderModel.Status,
		&orderModel.CreatedAt,
		&orderModel.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeOrderProto(orderModel), nil
}

// Convert db model to proto model
func makeOrderProto(order models.Order) *order_service.Order {
	return &order_service.Order{
		Id:                order.Id,
		ClientId:          order.ClientId,
		DeliveryLatitude:  order.DeliveryLatitude,
		DeliveryLongitude: order.DeliveryLongitude,
		TotalPrice:        order.TotalPrice,
		Status:            order.Status,
		CreatedAt:         timestamppb.New(order.CreatedAt),
		UpdatedAt:         timestamppb.New(order.UpdatedAt),
	}
}

// Convert proto model to db model
func makeOrderModel(order *order_service.Order) models.Order {
	return models.Order{
		Id:                order.Id,
		ClientId:          order.ClientId,
		DeliveryLatitude:  order.DeliveryLatitude,
		DeliveryLongitude: order.DeliveryLongitude,
		TotalPrice:        order.TotalPrice,
		Status:            order.Status,
		CreatedAt:         order.CreatedAt.AsTime(),
		UpdatedAt:         order.UpdatedAt.AsTime(),
	}
}
