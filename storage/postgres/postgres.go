package postgres

import (
	"context"
	"fmt"

	"github.com/flash_sale/flash_sale_order_service/config"
	"github.com/flash_sale/flash_sale_order_service/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// StoragePg implements the storage.StorageI interface for PostgreSQL.
type StoragePg struct {
	db             *pgx.Conn
	basketRepo     storage.BasketI
	basketItemRepo storage.BasketItemI
	orderRepo      storage.OrderI
	orderItemRepo  storage.OrderItemI
}

// NewStoragePg creates a new PostgreSQL storage instance.
func NewStoragePg(cfg config.Config) (storage.StorageI, error) {
	dbCon := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	db, err := pgx.Connect(context.Background(), dbCon)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}

	if err = db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error pinging postgres: %w", err)
	}

	return &StoragePg{
		db:             db,
		basketRepo:     NewBasketRepo(db),
		basketItemRepo: NewBasketItemRepo(db),
		orderRepo:      NewOrderRepo(db),
		orderItemRepo:  NewOrderItemRepo(db),
	}, nil
}

// Close closes the PostgreSQL connection.
func (s *StoragePg) Close() {
	if err := s.db.Close(context.Background()); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			fmt.Printf("Error closing database connection: %s (Code: %s)\n", pgErr.Message, pgErr.Code)
		} else {
			fmt.Printf("Error closing database connection: %s\n", err.Error())
		}
	}
}

// Basket returns the BasketI implementation for PostgreSQL.
func (s *StoragePg) Basket() storage.BasketI {
	return s.basketRepo
}

// BasketItem returns the BasketItemI implementation for PostgreSQL.
func (s *StoragePg) BasketItem() storage.BasketItemI {
	return s.basketItemRepo
}

// Order returns the OrderI implementation for PostgreSQL.
func (s *StoragePg) Order() storage.OrderI {
	return s.orderRepo
}

// OrderItem returns the OrderItemI implementation for PostgreSQL.
func (s *StoragePg) OrderItem() storage.OrderItemI {
	return s.orderItemRepo
}
