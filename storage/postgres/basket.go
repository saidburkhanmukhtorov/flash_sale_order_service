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

// ... (other code) ...

type BasketRepo struct {
	db *pgx.Conn
}

func NewBasketRepo(db *pgx.Conn) *BasketRepo {
	return &BasketRepo{
		db: db,
	}
}

func (r *BasketRepo) CreateBasket(ctx context.Context, req *order_service.CreateBasketRequest) (*order_service.Basket, error) {
	if req.Basket.Id == "" {
		req.Basket.Id = uuid.NewString()
	}

	query := `
		INSERT INTO baskets (
			id,
			user_id,
			status,
			created_at,
			updated_at,
			deleted_at
		) VALUES (
			$1, $2, $3, NOW(), NOW(), 0
		) RETURNING id, created_at, updated_at
	`

	basketModel := makeBasketModel(req.Basket)

	err := r.db.QueryRow(ctx, query,
		basketModel.Id,
		basketModel.UserId,
		basketModel.Status,
	).Scan(&basketModel.Id, &basketModel.CreatedAt, &basketModel.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return makeBasketProto(basketModel), nil
}

func (r *BasketRepo) GetBasket(ctx context.Context, req *order_service.GetBasketRequest) (*order_service.Basket, error) {
	var basketModel models.Basket

	query := `
		SELECT 
			id,
			user_id,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM baskets
		WHERE id = $1 AND deleted_at = 0
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&basketModel.Id,
		&basketModel.UserId,
		&basketModel.Status,
		&basketModel.CreatedAt,
		&basketModel.UpdatedAt,
		&basketModel.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeBasketProto(basketModel), nil
}

func (r *BasketRepo) UpdateBasket(ctx context.Context, req *order_service.UpdateBasketRequest) (*order_service.Basket, error) {
	query := `
		UPDATE baskets
		SET 
			user_id = $1,
			status = $2,
			updated_at = NOW()
		WHERE id = $3 AND deleted_at = 0
		RETURNING id, user_id, status, created_at, updated_at
	`

	basketModel := makeBasketModel(req.Basket)

	err := r.db.QueryRow(ctx, query,
		basketModel.UserId,
		basketModel.Status,
		basketModel.Id,
	).Scan(
		&basketModel.Id,
		&basketModel.UserId,
		&basketModel.Status,
		&basketModel.CreatedAt,
		&basketModel.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeBasketProto(basketModel), nil
}

func (r *BasketRepo) DeleteBasket(ctx context.Context, req *order_service.DeleteBasketRequest) (*order_service.DeleteBasketResponse, error) {
	query := `
		UPDATE baskets
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at = 0
	`

	_, err := r.db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		return nil, err
	}

	return &order_service.DeleteBasketResponse{
		Message: "Basket deleted successfully",
	}, nil
}

func (r *BasketRepo) ListBaskets(ctx context.Context, req *order_service.ListBasketsRequest) (*order_service.ListBasketsResponse, error) {
	var args []interface{}
	count := 1
	query := `
		SELECT 
			id,
			user_id,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM 
			baskets
		WHERE 1=1 AND deleted_at = 0
	`

	filter := ""

	if req.UserId != "" {
		filter += fmt.Sprintf(" AND user_id = $%d", count)
		args = append(args, req.UserId)
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

	totalCountQuery := "SELECT count(*) FROM baskets WHERE 1=1 AND deleted_at = 0" + filter
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

	var basketList []*order_service.Basket

	for rows.Next() {
		var basketModel models.Basket
		err = rows.Scan(
			&basketModel.Id,
			&basketModel.UserId,
			&basketModel.Status,
			&basketModel.CreatedAt,
			&basketModel.UpdatedAt,
			&basketModel.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		basketList = append(basketList, makeBasketProto(basketModel))
	}

	return &order_service.ListBasketsResponse{
		Baskets: basketList,
		Total:   int32(totalCount),
	}, nil
}

func (r *BasketRepo) UpdateBasketStatus(ctx context.Context, req *order_service.UpdateBasketStatusRequest) (*order_service.Basket, error) {
	query := `
		UPDATE baskets
		SET 
			status = $1,
			updated_at = NOW()
		WHERE id = $2 AND deleted_at = 0
		RETURNING id, user_id, status, created_at, updated_at
	`

	var basketModel models.Basket

	err := r.db.QueryRow(ctx, query,
		req.Status,
		req.Id,
	).Scan(
		&basketModel.Id,
		&basketModel.UserId,
		&basketModel.Status,
		&basketModel.CreatedAt,
		&basketModel.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return makeBasketProto(basketModel), nil
}

// Convert db model to proto model
func makeBasketProto(basket models.Basket) *order_service.Basket {
	return &order_service.Basket{
		Id:        basket.Id,
		UserId:    basket.UserId,
		Status:    basket.Status,
		CreatedAt: timestamppb.New(basket.CreatedAt),
		UpdatedAt: timestamppb.New(basket.UpdatedAt),
	}
}

// Convert proto model to db model
func makeBasketModel(basket *order_service.Basket) models.Basket {
	return models.Basket{
		Id:        basket.Id,
		UserId:    basket.UserId,
		Status:    basket.Status,
		CreatedAt: basket.CreatedAt.AsTime(),
		UpdatedAt: basket.UpdatedAt.AsTime(),
	}
}
