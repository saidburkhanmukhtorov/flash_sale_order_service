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

// BasketToOrderConsumer consumes Kafka messages for converting basket items to order items.
type BasketToOrderConsumer struct {
	reader  *kafka.Reader
	storage storage.StorageI
}

// NewBasketToOrderConsumer creates a new BasketToOrderConsumer instance.
func NewBasketToOrderConsumer(kafkaBrokers []string, topic string, storage storage.StorageI) *BasketToOrderConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaBrokers,
		Topic:   topic,
		GroupID: "basket-to-order-group", // Choose a suitable group ID
	})
	return &BasketToOrderConsumer{reader: reader, storage: storage}
}

// Consume starts consuming messages from the Kafka topic.
func (c *BasketToOrderConsumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("error fetching message: %w", err)
		}

		switch string(msg.Key) {
		case "basket.convert_to_order":
			var convertModel order_service.ConvertBasketToOrderItemsRequest
			if err := json.Unmarshal(msg.Value, &convertModel); err != nil {
				log.Printf("error unmarshalling convert basket to order message: %v", err)
				continue
			}

			// Convert basket items to order items
			if _, err := c.storage.OrderItem().ConvertBasketToOrderItems(ctx, &convertModel); err != nil {
				log.Printf("error converting basket to order: %v", err)
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
