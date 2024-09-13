package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/flash_sale/flash_sale_order_service/config"
	consumer "github.com/flash_sale/flash_sale_order_service/kafka"

	"github.com/flash_sale/flash_sale_order_service/genproto/order_service"
	"github.com/flash_sale/flash_sale_order_service/service"
	"github.com/flash_sale/flash_sale_order_service/storage/postgres"
	"github.com/flash_sale/flash_sale_order_service/storage/redis"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	// Initialize PostgreSQL storage
	pgStorage, err := postgres.NewStoragePg(cfg)
	if err != nil {
		log.Fatalf("failed to initialize PostgreSQL storage: %v", err)
	}

	// Initialize Redis client
	redisClient, err := redis.Connect(&cfg)
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize Kafka consumers
	basketItemConsumer := consumer.NewBasketItemConsumer(
		cfg.KafkaBrokers,
		"basket_item_topic",
		pgStorage,
	)
	basketToOrderConsumer := consumer.NewBasketToOrderConsumer(
		cfg.KafkaBrokers,
		"basket_to_order_topic",
		pgStorage,
	)

	// Start consumers in separate goroutines
	go func() {
		log.Println("basket_item_topic is ready to accept requests.")
		if err := basketItemConsumer.Consume(context.Background()); err != nil {
			log.Fatalf("basket item consumer error: %v", err)
		}
	}()

	go func() {
		log.Println("basket_to_order_topic is ready to accept requests.")
		if err := basketToOrderConsumer.Consume(context.Background()); err != nil {
			log.Fatalf("basket to order consumer error: %v", err)
		}
	}()

	// Initialize gRPC server
	lis, err := net.Listen("tcp", cfg.OrderServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Register gRPC services
	order_service.RegisterBasketServiceServer(s, service.NewBasketService(pgStorage, redisClient))
	order_service.RegisterBasketItemServiceServer(s, service.NewBasketItemService(pgStorage))
	order_service.RegisterOrderServiceServer(s, service.NewOrderService(pgStorage, redisClient))
	order_service.RegisterOrderItemServiceServer(s, service.NewOrderItemService(pgStorage, redisClient))

	fmt.Printf("server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
