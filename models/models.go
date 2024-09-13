package models

import "time"

// Basket represents a shopping basket model for the database.
type Basket struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	Status    string    `db:"status"` // Possible values: 'OPEN', 'CHECKED_OUT'
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt int64     `db:"deleted_at"`
}

// BasketItem represents a basket item model for the database.
type BasketItem struct {
	Id                      string    `db:"id"`
	BasketId                string    `db:"basket_id"`
	ProductId               string    `db:"product_id"`
	FlashSaleEventProductId string    `db:"flash_sale_event_product_id"`
	DiscountProductId       string    `db:"discount_product_id"`
	Quantity                int32     `db:"quantity"`
	UnitPrice               float32   `db:"unit_price"`
	TotalPrice              float32   `db:"total_price"`
	ProductType             string    `db:"product_type"` // Possible values: 'REGULAR', 'FLASH_SALE', 'DISCOUNT'
	CreatedAt               time.Time `db:"created_at"`
	UpdatedAt               time.Time `db:"updated_at"`
	DeletedAt               int64     `db:"deleted_at"`
}

// Order represents an order model for the database.
type Order struct {
	Id                string    `db:"id"`
	ClientId          string    `db:"client_id"`
	DeliveryLatitude  float64   `db:"delivery_latitude"`
	DeliveryLongitude float64   `db:"delivery_longitude"`
	TotalPrice        float32   `db:"total_price"`
	Status            string    `db:"status"` // Possible values: 'PENDING', 'PROCESSING', 'SHIPPED', 'DELIVERED', 'CANCELLED'
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	DeletedAt         int64     `db:"deleted_at"`
}

// OrderItem represents an order item model for the database.
type OrderItem struct {
	Id                      string    `db:"id"`
	OrderId                 string    `db:"order_id"`
	ProductId               string    `db:"product_id"`
	FlashSaleEventProductId string    `db:"flash_sale_event_product_id"`
	DiscountProductId       string    `db:"discount_product_id"`
	Quantity                int32     `db:"quantity"`
	UnitPrice               float32   `db:"unit_price"`
	TotalPrice              float32   `db:"total_price"`
	DiscountApplied         float32   `db:"discount_applied"`
	ProductType             string    `db:"product_type"` // Possible values: 'REGULAR', 'FLASH_SALE', 'DISCOUNT'
	CreatedAt               time.Time `db:"created_at"`
	UpdatedAt               time.Time `db:"updated_at"`
	DeletedAt               int64     `db:"deleted_at"`
}

// Database Model
type Product struct {
	Id            string    `db:"id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	BasePrice     float32   `db:"base_price"`
	CurrentPrice  float32   `db:"current_price"`
	ImageUrl      string    `db:"image_url"`
	StockQuantity int32     `db:"stock_quantity"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	DeletedAt     int64     `db:"deleted_at"`
}

// Discount represents a discount model for the database.
type Discount struct {
	Id            string    `db:"id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	DiscountType  string    `db:"discount_type"` // Possible values: 'PERCENTAGE', 'FIXED_AMOUNT'
	DiscountValue float32   `db:"discount_value"`
	StartDate     time.Time `db:"start_date"`
	EndDate       time.Time `db:"end_date"`
	IsActive      bool      `db:"is_active"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	DeletedAt     int64     `db:"deleted_at"`
}

// FlashSaleEvent represents a flash sale event model for the database.
type FlashSaleEvent struct {
	Id          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	StartTime   time.Time `db:"start_time"`
	EndTime     time.Time `db:"end_time"`
	Status      string    `db:"status"`     // Possible values: 'UPCOMING', 'ACTIVE', 'ENDED'
	EventType   string    `db:"event_type"` // Possible values: 'FLASH_SALE', 'PROMOTION'
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	DeletedAt   int64     `db:"deleted_at"`
}

// ProductDiscount represents a product discount model for the database.
type ProductDiscount struct {
	Id         string    `db:"id"`
	ProductId  string    `db:"product_id"`
	DiscountId string    `db:"discount_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	DeletedAt  int64     `db:"deleted_at"`
}

// FlashSaleEventProduct represents a flash sale event product model for the database.
type FlashSaleEventProduct struct {
	Id                 string    `db:"id"`
	EventId            string    `db:"event_id"`
	ProductId          string    `db:"product_id"`
	DiscountPercentage float32   `db:"discount_percentage"`
	SalePrice          float32   `db:"sale_price"`
	AvailableQuantity  int32     `db:"available_quantity"`
	OriginalStock      int32     `db:"original_stock"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
	DeletedAt          int64     `db:"deleted_at"`
}
