package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/flash_sale/flash_sale_order_service/config"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// Client represents a Redis client.
type Client struct {
	*redis.Client
}

// Connect establishes a connection to the Redis server.
func Connect(cfg *config.Config) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test the connection
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &Client{client}, nil
}

// Add these new methods for notification handling
type Notification struct {
	ID      string    `json:"id"`
	UserID  string    `json:"user_id"`
	Message string    `json:"message"`
	Created time.Time `json:"created"`
}

func (c *Client) AddNotification(ctx context.Context, user_id, message string) error {
	notification := Notification{
		ID:      uuid.NewString(),
		UserID:  user_id,
		Message: message,
		Created: time.Now(),
	}
	json, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	pipe := c.Pipeline()

	// Add to unread set
	pipe.SAdd(ctx, fmt.Sprintf("unread:%s", notification.UserID), notification.ID)

	// Add to sorted set for ordering
	pipe.ZAdd(ctx, fmt.Sprintf("notifications:%s", notification.UserID), &redis.Z{
		Score:  float64(notification.Created.Unix()),
		Member: notification.ID,
	})

	// Store notification data
	pipe.Set(ctx, fmt.Sprintf("notification:%s", notification.ID), json, 0)

	_, err = pipe.Exec(ctx)
	return err
}
