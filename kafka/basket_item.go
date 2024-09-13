package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/flash_sale/flash_sale_order_service/genproto/order_service"
	"github.com/flash_sale/flash_sale_order_service/storage"
	"github.com/segmentio/kafka-go"
)

// BasketItemConsumer consumes Kafka messages related to basket items.
type BasketItemConsumer struct {
	reader  *kafka.Reader
	storage storage.StorageI
}

// NewBasketItemConsumer creates a new BasketItemConsumer instance.
func NewBasketItemConsumer(kafkaBrokers []string, topic string, storage storage.StorageI) *BasketItemConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   topic,
		GroupID: "basket-item-group", // Choose a suitable group ID
	})
	return &BasketItemConsumer{reader: reader, storage: storage}
}

// Consume starts consuming messages from the Kafka topic.
func (c *BasketItemConsumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("error fetching message: %w", err)
		}

		switch string(msg.Key) {
		case "basket_item.create":
			var createModel order_service.CreateBasketItemRequest
			if err := json.Unmarshal(msg.Value, &createModel); err != nil {
				log.Printf("error unmarshalling create basket item message: %v", err)
				continue
			}

			// Create the basket item in the database
			if _, err := c.storage.BasketItem().CreateBasketItem(ctx, &createModel); err != nil {
				log.Printf("error creating basket item: %v", err)
				continue
			}

		default:
			log.Printf("unknown message key: %s", msg.Key)
		}

		// Commit the message
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("error committing message: %w", err)
		}
	}
}
