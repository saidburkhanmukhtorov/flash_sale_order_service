package test

import (
	"context"
	"testing"
	"time"

	"github.com/flash_sale/flash_sale_order_service/genproto/order_service"
	"github.com/flash_sale/flash_sale_order_service/storage/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepo(t *testing.T) {
	db := createDBConnection(t) // Use the existing createDBConnection function
	defer db.Close(context.Background())

	// Initialize repositories
	basketRepo := postgres.NewBasketRepo(db)
	basketItemRepo := postgres.NewBasketItemRepo(db)
	orderRepo := postgres.NewOrderRepo(db)
	orderItemRepo := postgres.NewOrderItemRepo(db)

	// 1. Create a user
	userID := uuid.NewString()
	createUser(t, db, userID)
	defer deleteUser(t, db, userID)

	// 2. Create products
	product1ID := uuid.NewString()
	createProduct(t, db, product1ID, "Product 1", 10.0)
	defer deleteProduct(t, db, product1ID)

	product2ID := uuid.NewString()
	createProduct(t, db, product2ID, "Product 2", 20.0)
	defer deleteProduct(t, db, product2ID)

	// 3. Create a flash sale event
	flashSaleEventID := uuid.NewString()
	createFlashSaleEvent(t, db, flashSaleEventID, "Flash Sale Event", time.Now().Add(time.Hour), time.Now().Add(24*time.Hour), "ACTIVE")
	defer deleteFlashSaleEvent(t, db, flashSaleEventID)

	// 4. Create a discount
	discountID := uuid.NewString()
	createDiscount(t, db, discountID, "Discount 1", "PERCENTAGE", 10.0, true)
	defer deleteDiscount(t, db, discountID)

	// 5. Add product to the discount
	productDiscountID := uuid.NewString()
	createProductDiscount(t, db, productDiscountID, product1ID, discountID)
	defer deleteProductDiscount(t, db, productDiscountID)

	// 6. Add product to the flash sale event
	flashSaleEventProductID := uuid.NewString()
	createFlashSaleEventProduct(t, db, flashSaleEventProductID, flashSaleEventID, product2ID, 20.0, 18.0)
	defer deleteFlashSaleEventProduct(t, db, flashSaleEventProductID)

	// --- Basket Tests ---

	t.Run("CreateBasket", func(t *testing.T) {
		basket, err := basketRepo.CreateBasket(context.Background(), &order_service.CreateBasketRequest{
			Basket: &order_service.Basket{
				UserId: userID,
				Status: "OPEN",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, basket)
		assert.NotEmpty(t, basket.Id)
		defer deleteBasket(t, db, basket.Id)
	})

	t.Run("GetBasket", func(t *testing.T) {
		// Create a basket first
		createdBasket, err := basketRepo.CreateBasket(context.Background(), &order_service.CreateBasketRequest{
			Basket: &order_service.Basket{
				UserId: userID,
				Status: "OPEN",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdBasket)

		basket, err := basketRepo.GetBasket(context.Background(), &order_service.GetBasketRequest{Id: createdBasket.Id})
		assert.NoError(t, err)
		assert.NotNil(t, basket)
		assert.Equal(t, createdBasket.Id, basket.Id)

		defer deleteBasket(t, db, createdBasket.Id)
	})

	t.Run("UpdateBasket", func(t *testing.T) {
		// Create a basket first
		createdBasket, err := basketRepo.CreateBasket(context.Background(), &order_service.CreateBasketRequest{
			Basket: &order_service.Basket{
				UserId: userID,
				Status: "OPEN",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdBasket)

		// Update the basket
		createdBasket.Status = "CHECKED_OUT"
		updatedBasket, err := basketRepo.UpdateBasket(context.Background(), &order_service.UpdateBasketRequest{
			Basket: createdBasket,
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedBasket)
		assert.Equal(t, "CHECKED_OUT", updatedBasket.Status)

		defer deleteBasket(t, db, createdBasket.Id)
	})

	t.Run("DeleteBasket", func(t *testing.T) {
		// Create a basket first
		createdBasket, err := basketRepo.CreateBasket(context.Background(), &order_service.CreateBasketRequest{
			Basket: &order_service.Basket{
				UserId: userID,
				Status: "OPEN",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdBasket)

		// Delete the basket
		_, err = basketRepo.DeleteBasket(context.Background(), &order_service.DeleteBasketRequest{Id: createdBasket.Id})
		assert.NoError(t, err)

		// Try to get the deleted basket
		_, err = basketRepo.GetBasket(context.Background(), &order_service.GetBasketRequest{Id: createdBasket.Id})
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	})

	t.Run("ListBaskets", func(t *testing.T) {
		// Create a few test baskets
		basketsToCreate := []*order_service.Basket{
			{
				UserId: userID,
				Status: "OPEN",
			},
			{
				UserId: userID,
				Status: "CHECKED_OUT",
			},
		}

		for _, basket := range basketsToCreate {
			_, err := basketRepo.CreateBasket(context.Background(), &order_service.CreateBasketRequest{Basket: basket})
			assert.NoError(t, err)
			defer deleteBasket(t, db, basket.Id)
		}

		// Test ListBaskets with user ID filter
		baskets, err := basketRepo.ListBaskets(context.Background(), &order_service.ListBasketsRequest{
			UserId: userID,
			Page:   1,
			Limit:  10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, baskets)
		assert.GreaterOrEqual(t, len(baskets.Baskets), 2)
	})

	t.Run("UpdateBasketStatus", func(t *testing.T) {
		// Create a basket first
		createdBasket, err := basketRepo.CreateBasket(context.Background(), &order_service.CreateBasketRequest{
			Basket: &order_service.Basket{
				UserId: userID,
				Status: "OPEN",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdBasket)

		// Update the basket status
		updatedBasket, err := basketRepo.UpdateBasketStatus(context.Background(), &order_service.UpdateBasketStatusRequest{
			Id:     createdBasket.Id,
			Status: "CHECKED_OUT",
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedBasket)
		assert.Equal(t, "CHECKED_OUT", updatedBasket.Status)

		defer deleteBasket(t, db, createdBasket.Id)
	})

	// --- Basket Item Tests ---

	t.Run("CreateBasketItem", func(t *testing.T) {
		// Create a basket first
		createdBasket, err := basketRepo.CreateBasket(context.Background(), &order_service.CreateBasketRequest{
			Basket: &order_service.Basket{
				UserId: userID,
				Status: "OPEN",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdBasket)

		// Create a basket item
		basketItem, err := basketItemRepo.CreateBasketItem(context.Background(), &order_service.CreateBasketItemRequest{
			BasketItem: &order_service.BasketItem{
				BasketId:                createdBasket.Id,
				ProductId:               product1ID,
				FlashSaleEventProductId: "", // Not a flash sale item
				DiscountProductId:       "", // Not a discount item
				Quantity:                2,
				UnitPrice:               10.0,
				TotalPrice:              20.0,
				ProductType:             "REGULAR",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, basketItem)
		assert.NotEmpty(t, basketItem.Id)

		defer deleteBasketItem(t, db, basketItem.Id)
		defer deleteBasket(t, db, createdBasket.Id)
	})

	t.Run("GetBasketItem", func(t *testing.T) {
		// Create a basket first
		createdBasket, err := basketRepo.CreateBasket(context.Background(), &order_service.CreateBasketRequest{
			Basket: &order_service.Basket{
				UserId: userID,
				Status: "OPEN",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdBasket)

		// Create a basket item
		createdBasketItem, err := basketItemRepo.CreateBasketItem(context.Background(), &order_service.CreateBasketItemRequest{
			BasketItem: &order_service.BasketItem{
				BasketId:                createdBasket.Id,
				ProductId:               product1ID,
				FlashSaleEventProductId: "", // Not a flash sale item
				DiscountProductId:       "", // Not a discount item
				Quantity:                2,
				UnitPrice:               10.0,
				TotalPrice:              20.0,
				ProductType:             "REGULAR",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdBasketItem)

		// Get the basket item
		basketItem, err := basketItemRepo.GetBasketItem(context.Background(), &order_service.GetBasketItemRequest{Id: createdBasketItem.Id})
		assert.NoError(t, err)
		assert.NotNil(t, basketItem)
		assert.Equal(t, createdBasketItem.Id, basketItem.Id)

		defer deleteBasketItem(t, db, createdBasketItem.Id)
		defer deleteBasket(t, db, createdBasket.Id)
	})

	t.Run("DeleteBasketItem", func(t *testing.T) {
		// Create a basket first
		createdBasket, err := basketRepo.CreateBasket(context.Background(), &order_service.CreateBasketRequest{
			Basket: &order_service.Basket{
				UserId: userID,
				Status: "OPEN",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdBasket)

		// Create a basket item
		createdBasketItem, err := basketItemRepo.CreateBasketItem(context.Background(), &order_service.CreateBasketItemRequest{
			BasketItem: &order_service.BasketItem{
				BasketId:                createdBasket.Id,
				ProductId:               product1ID,
				FlashSaleEventProductId: "", // Not a flash sale item
				DiscountProductId:       "", // Not a discount item
				Quantity:                2,
				UnitPrice:               10.0,
				TotalPrice:              20.0,
				ProductType:             "REGULAR",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdBasketItem)

		// Delete the basket item
		_, err = basketItemRepo.DeleteBasketItem(context.Background(), &order_service.DeleteBasketItemRequest{Id: createdBasketItem.Id})
		assert.NoError(t, err)

		// Try to get the deleted basket item
		_, err = basketItemRepo.GetBasketItem(context.Background(), &order_service.GetBasketItemRequest{Id: createdBasketItem.Id})
		assert.ErrorIs(t, err, pgx.ErrNoRows)

		defer deleteBasket(t, db, createdBasket.Id)
	})

	t.Run("ListBasketItems", func(t *testing.T) {
		// Create a basket first
		createdBasket, err := basketRepo.CreateBasket(context.Background(), &order_service.CreateBasketRequest{
			Basket: &order_service.Basket{
				UserId: userID,
				Status: "OPEN",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdBasket)

		// Create a few test basket items
		basketItemsToCreate := []*order_service.BasketItem{
			{
				BasketId:                createdBasket.Id,
				ProductId:               product1ID,
				FlashSaleEventProductId: "", // Not a flash sale item
				DiscountProductId:       "", // Not a discount item
				Quantity:                2,
				UnitPrice:               10.0,
				TotalPrice:              20.0,
				ProductType:             "REGULAR",
			},
			{
				BasketId:                createdBasket.Id,
				ProductId:               product2ID,
				FlashSaleEventProductId: flashSaleEventProductID,
				DiscountProductId:       "",
				Quantity:                1,
				UnitPrice:               18.0,
				TotalPrice:              18.0,
				ProductType:             "FLASH_SALE",
			},
		}

		for _, basketItem := range basketItemsToCreate {
			_, err := basketItemRepo.CreateBasketItem(context.Background(), &order_service.CreateBasketItemRequest{BasketItem: basketItem})
			assert.NoError(t, err)
			defer deleteBasketItem(t, db, basketItem.Id)
		}

		// Test ListBasketItems with basket ID filter
		basketItems, err := basketItemRepo.ListBasketItems(context.Background(), &order_service.ListBasketItemsRequest{
			BasketId: createdBasket.Id,
			Page:     1,
			Limit:    10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, basketItems)
		assert.GreaterOrEqual(t, len(basketItems.BasketItems), 2)

		defer deleteBasket(t, db, createdBasket.Id)
	})

	// --- Order Tests ---

	t.Run("CreateOrder", func(t *testing.T) {
		order, err := orderRepo.CreateOrder(context.Background(), &order_service.CreateOrderRequest{
			Order: &order_service.Order{
				ClientId:          userID,
				DeliveryLatitude:  37.7749,
				DeliveryLongitude: -122.4194,
				TotalPrice:        0.0, // Will be updated later
				Status:            "PENDING",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.NotEmpty(t, order.Id)
		defer deleteOrder(t, db, order.Id)
	})

	t.Run("GetOrder", func(t *testing.T) {
		// Create an order first
		createdOrder, err := orderRepo.CreateOrder(context.Background(), &order_service.CreateOrderRequest{
			Order: &order_service.Order{
				ClientId:          userID,
				DeliveryLatitude:  37.7749,
				DeliveryLongitude: -122.4194,
				TotalPrice:        0.0, // Will be updated later
				Status:            "PENDING",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdOrder)

		order, err := orderRepo.GetOrder(context.Background(), &order_service.GetOrderRequest{Id: createdOrder.Id})
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, createdOrder.Id, order.Id)

		defer deleteOrder(t, db, createdOrder.Id)
	})

	t.Run("UpdateOrder", func(t *testing.T) {
		// Create an order first
		createdOrder, err := orderRepo.CreateOrder(context.Background(), &order_service.CreateOrderRequest{
			Order: &order_service.Order{
				ClientId:          userID,
				DeliveryLatitude:  37.7749,
				DeliveryLongitude: -122.4194,
				TotalPrice:        0.0, // Will be updated later
				Status:            "PENDING",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdOrder)

		// Update the order
		createdOrder.Status = "PROCESSING"
		updatedOrder, err := orderRepo.UpdateOrder(context.Background(), &order_service.UpdateOrderRequest{
			Order: createdOrder,
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedOrder)
		assert.Equal(t, "PROCESSING", updatedOrder.Status)

		defer deleteOrder(t, db, createdOrder.Id)
	})

	t.Run("DeleteOrder", func(t *testing.T) {
		// Create an order first
		createdOrder, err := orderRepo.CreateOrder(context.Background(), &order_service.CreateOrderRequest{
			Order: &order_service.Order{
				ClientId:          userID,
				DeliveryLatitude:  37.7749,
				DeliveryLongitude: -122.4194,
				TotalPrice:        0.0, // Will be updated later
				Status:            "PENDING",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdOrder)

		// Delete the order
		_, err = orderRepo.DeleteOrder(context.Background(), &order_service.DeleteOrderRequest{Id: createdOrder.Id})
		assert.NoError(t, err)

		// Try to get the deleted order
		_, err = orderRepo.GetOrder(context.Background(), &order_service.GetOrderRequest{Id: createdOrder.Id})
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	})

	t.Run("ListOrders", func(t *testing.T) {
		// Create a few test orders
		ordersToCreate := []*order_service.Order{
			{
				ClientId:          userID,
				DeliveryLatitude:  37.7749,
				DeliveryLongitude: -122.4194,
				TotalPrice:        10.0,
				Status:            "PENDING",
			},
			{
				ClientId:          userID,
				DeliveryLatitude:  34.0522,
				DeliveryLongitude: -118.2437,
				TotalPrice:        20.0,
				Status:            "PROCESSING",
			},
		}

		for _, order := range ordersToCreate {
			_, err := orderRepo.CreateOrder(context.Background(), &order_service.CreateOrderRequest{Order: order})
			assert.NoError(t, err)
			defer deleteOrder(t, db, order.Id)
		}

		// Test ListOrders with client ID and status filters
		orders, err := orderRepo.ListOrders(context.Background(), &order_service.ListOrdersRequest{
			ClientId: userID,
			Status:   "PENDING",
			Page:     1,
			Limit:    10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, orders)
		assert.GreaterOrEqual(t, len(orders.Orders), 1)
	})

	t.Run("UpdateOrderStatus", func(t *testing.T) {
		// Create an order first
		createdOrder, err := orderRepo.CreateOrder(context.Background(), &order_service.CreateOrderRequest{
			Order: &order_service.Order{
				ClientId:          userID,
				DeliveryLatitude:  37.7749,
				DeliveryLongitude: -122.4194,
				TotalPrice:        0.0, // Will be updated later
				Status:            "PENDING",
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, createdOrder)

		// Update the order status
		updatedOrder, err := orderRepo.UpdateOrderStatus(context.Background(), &order_service.UpdateOrderStatusRequest{
			Id:     createdOrder.Id,
			Status: "SHIPPED",
		})
		assert.NoError(t, err)
		assert.NotNil(t, updatedOrder)
		assert.Equal(t, "SHIPPED", updatedOrder.Status)

		defer deleteOrder(t, db, createdOrder.Id)
	})

	// --- Order Item Tests ---

	t.Run("ConvertBasketToOrderItems", func(t *testing.T) {
		// Create a basket
		basketID := uuid.NewString()
		createBasket(t, db, basketID, userID, "OPEN")
		defer deleteBasket(t, db, basketID)

		// Create basket items
		basketItem1ID := uuid.NewString()
		createBasketItemRegular(t, db, basketItem1ID, basketID, product1ID, 2, 10.0, 20.0)
		defer deleteBasketItem(t, db, basketItem1ID)

		basketItem2ID := uuid.NewString()
		createBasketItemFlashSale(t, db, basketItem2ID, basketID, product2ID, flashSaleEventProductID, 1, 18.0, 18.0)
		defer deleteBasketItem(t, db, basketItem2ID)

		// Create an order
		orderID := uuid.NewString()
		createOrder(t, db, orderID, userID, 0, 0, 0, "PENDING")
		defer deleteOrder(t, db, orderID)

		// Convert basket items to order items
		returnedOrderID, err := orderItemRepo.ConvertBasketToOrderItems(context.Background(), &order_service.ConvertBasketToOrderItemsRequest{
			BasketId: basketID,
			OrderId:  orderID,
		})
		assert.NoError(t, err)
		assert.Equal(t, orderID, returnedOrderID.Id)

		// Check if order items were created
		orderItems, err := orderItemRepo.ListOrderItems(context.Background(), &order_service.ListOrderItemsRequest{
			OrderId: orderID,
			Page:    1,
			Limit:   10,
		})
		assert.NoError(t, err)
		assert.NotNil(t, orderItems)
		assert.Equal(t, 2, len(orderItems.OrderItems))

		// Check if order total price is updated
		order, err := orderRepo.GetOrder(context.Background(), &order_service.GetOrderRequest{Id: orderID})
		assert.NoError(t, err)
		assert.Equal(t, float32(38.0), order.TotalPrice) // 20.0 (regular) + 18.0 (flash sale)
	})
	t.Run("DeleteOrderItem", func(t *testing.T) {
		// Create a basket
		basketID := uuid.NewString()
		createBasket(t, db, basketID, userID, "OPEN")
		defer deleteBasket(t, db, basketID)

		// Create a basket item
		basketItem1ID := uuid.NewString()
		createBasketItemRegular(t, db, basketItem1ID, basketID, product1ID, 2, 10.0, 20.0)
		defer deleteBasketItem(t, db, basketItem1ID)

		// Create an order
		orderID := uuid.NewString()
		createOrder(t, db, orderID, userID, 0, 0, 20.0, "PENDING") // Initial total price is 20.0
		defer deleteOrder(t, db, orderID)

		// Convert basket items to order items
		_, err := orderItemRepo.ConvertBasketToOrderItems(context.Background(), &order_service.ConvertBasketToOrderItemsRequest{
			BasketId: basketID,
			OrderId:  orderID,
		})
		assert.NoError(t, err)

		// Get the created order item ID
		orderItems, err := orderItemRepo.ListOrderItems(context.Background(), &order_service.ListOrderItemsRequest{
			OrderId: orderID,
			Page:    1,
			Limit:   10,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(orderItems.OrderItems), 1) // At least one order item should exist
		orderItemID := orderItems.OrderItems[0].Id

		// Delete the order item using the orderItemID
		deletedOrderID, err := orderItemRepo.DeleteOrderItem(context.Background(), &order_service.DeleteOrderItemRequest{Id: orderItemID})
		assert.NoError(t, err)
		assert.Equal(t, orderID, deletedOrderID)

		// Verify that the order item is deleted
		_, err = orderItemRepo.GetOrderItem(context.Background(), &order_service.GetOrderItemRequest{Id: orderItemID})
		assert.ErrorIs(t, err, pgx.ErrNoRows)

		// Check if order total price is updated
		order, err := orderRepo.GetOrder(context.Background(), &order_service.GetOrderRequest{Id: orderID})
		assert.NoError(t, err)
		assert.Equal(t, float32(0.0), order.TotalPrice) // Total price should be 0 after deleting the only item
	})
}

// Helper functions to create and delete test data
func createUser(t *testing.T, db *pgx.Conn, userID string) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO users (id, username, email, password_hash, full_name, date_of_birth, role, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), 0)
	`, userID, uuid.NewString(), uuid.NewString()+"@example.com", "password", "Test User", "2000-01-01", "user")
	assert.NoError(t, err)
}

func deleteUser(t *testing.T, db *pgx.Conn, userID string) {
	// _, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID)
	// assert.NoError(t, err)
}

func createProduct(t *testing.T, db *pgx.Conn, productID, name string, price float32) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO products (id, name, description, base_price, current_price, image_url, stock_quantity, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), 0)
	`, productID, name, "Test Description", price, price, "https://example.com/image.jpg", 100)
	assert.NoError(t, err)
}

func deleteProduct(t *testing.T, db *pgx.Conn, productID string) {
	// _, err := db.Exec(context.Background(), "DELETE FROM products WHERE id = $1", productID)
	// assert.NoError(t, err)
}

func createFlashSaleEvent(t *testing.T, db *pgx.Conn, eventID, name string, startTime, endTime time.Time, status string) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO flash_sale_events (id, name, description, start_time, end_time, status, event_type, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), 0)
	`, eventID, name, "Test Description", startTime, endTime, status, "FLASH_SALE")
	assert.NoError(t, err)
}

func deleteFlashSaleEvent(t *testing.T, db *pgx.Conn, eventID string) {
	// _, err := db.Exec(context.Background(), "DELETE FROM flash_sale_events WHERE id = $1", eventID)
	// assert.NoError(t, err)
}

func createDiscount(t *testing.T, db *pgx.Conn, discountID, name, discountType string, discountValue float32, isActive bool) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO discounts (id, name, description, discount_type, discount_value, start_date, end_date, is_active, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), 0)
	`, discountID, name, "Test Description", discountType, discountValue, time.Now(), time.Now().Add(24*time.Hour), isActive)
	assert.NoError(t, err)
}

func deleteDiscount(t *testing.T, db *pgx.Conn, discountID string) {
	_, err := db.Exec(context.Background(), "DELETE FROM discounts WHERE id = $1", discountID)
	assert.NoError(t, err)
}

func createProductDiscount(t *testing.T, db *pgx.Conn, productDiscountID, productID, discountID string) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO product_discounts (id, product_id, discount_id, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, NOW(), NOW(), 0)
	`, productDiscountID, productID, discountID)
	assert.NoError(t, err)
}

func deleteProductDiscount(t *testing.T, db *pgx.Conn, productDiscountID string) {
	_, err := db.Exec(context.Background(), "DELETE FROM product_discounts WHERE id = $1", productDiscountID)
	assert.NoError(t, err)
}

func createFlashSaleEventProduct(t *testing.T, db *pgx.Conn, flashSaleEventProductID, eventID, productID string, discountPercentage, salePrice float32) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO flash_sale_event_products (id, event_id, product_id, discount_percentage, sale_price, available_quantity, original_stock, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), 0)
	`, flashSaleEventProductID, eventID, productID, discountPercentage, salePrice, 10, 10)
	assert.NoError(t, err)
}

func deleteFlashSaleEventProduct(t *testing.T, db *pgx.Conn, flashSaleEventProductID string) {
	// _, err := db.Exec(context.Background(), "DELETE FROM flash_sale_event_products WHERE id = $1", flashSaleEventProductID)
	// assert.NoError(t, err)
}

func createBasket(t *testing.T, db *pgx.Conn, basketID, userID, status string) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO baskets (id, user_id, status, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, NOW(), NOW(), 0)
	`, basketID, userID, status)
	assert.NoError(t, err)
}

func deleteBasket(t *testing.T, db *pgx.Conn, basketID string) {
	// _, err := db.Exec(context.Background(), "DELETE FROM baskets WHERE id = $1", basketID)
	// assert.NoError(t, err)
}

func createBasketItemRegular(t *testing.T, db *pgx.Conn, basketItemID, basketID, productID string, quantity int32, unitPrice, totalPrice float32) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO basket_items (id, basket_id, product_id, quantity, unit_price, total_price, product_type, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW(), 0)
	`, basketItemID, basketID, productID, quantity, unitPrice, totalPrice, "REGULAR")
	assert.NoError(t, err)
}

func createBasketItemFlashSale(t *testing.T, db *pgx.Conn, basketItemID, basketID, productID, flashSaleEventProductID string, quantity int32, unitPrice, totalPrice float32) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO basket_items (id, basket_id, product_id, flash_sale_event_product_id, quantity, unit_price, total_price, product_type, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), 0)
	`, basketItemID, basketID, productID, flashSaleEventProductID, quantity, unitPrice, totalPrice, "FLASH_SALE")
	assert.NoError(t, err)
}

func createBasketItemDiscount(t *testing.T, db *pgx.Conn, basketItemID, basketID, productID, discountProductID string, quantity int32, unitPrice, totalPrice float32) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO basket_items (id, basket_id, product_id, discount_product_id, quantity, unit_price, total_price, product_type, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), 0)
	`, basketItemID, basketID, productID, discountProductID, quantity, unitPrice, totalPrice, "DISCOUNT")
	assert.NoError(t, err)
}

func deleteBasketItem(t *testing.T, db *pgx.Conn, basketItemID string) {
	_, err := db.Exec(context.Background(), "DELETE FROM basket_items WHERE id = $1", basketItemID)
	assert.NoError(t, err)
}

func createOrder(t *testing.T, db *pgx.Conn, orderID, clientID string, deliveryLatitude, deliveryLongitude float64, totalPrice float32, status string) {
	_, err := db.Exec(context.Background(), `
		INSERT INTO orders (id, client_id, delivery_latitude, delivery_longitude, total_price, status, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), 0)
	`, orderID, clientID, deliveryLatitude, deliveryLongitude, totalPrice, status)
	assert.NoError(t, err)
}

func deleteOrder(t *testing.T, db *pgx.Conn, orderID string) {
	// _, err := db.Exec(context.Background(), "DELETE FROM orders WHERE id = $1", orderID)
	// assert.NoError(t, err)
}
