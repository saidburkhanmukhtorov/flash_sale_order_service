package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

// Config struct holds the configuration settings.
type Config struct {
	OrderServicePort string

	// PostgreSQL Configuration
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	KafkaBrokers     []string
	LOG_PATH         string

	// Redis Configuration
	RedisAddress  string
	RedisPassword string
	RedisDB       int
}

// Load loads the configuration from environment variables.
func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.OrderServicePort = cast.ToString(coalesce("ORDER_SERVICE_PORT", ":9091"))

	// PostgreSQL Configuration
	config.PostgresHost = cast.ToString(coalesce("POSTGRES_HOST", "postgres_dock"))
	config.PostgresPort = cast.ToInt(coalesce("POSTGRES_PORT", 5432))
	config.PostgresUser = cast.ToString(coalesce("POSTGRES_USER", "postgres"))
	config.PostgresPassword = cast.ToString(coalesce("POSTGRES_PASSWORD", "example"))
	config.PostgresDB = cast.ToString(coalesce("POSTGRES_DB", "memory"))

	// Redis Configuration
	config.RedisAddress = cast.ToString(coalesce("REDIS_ADDRESS", "redis:6379"))
	config.RedisPassword = cast.ToString(coalesce("REDIS_PASSWORD", ""))
	config.RedisDB = cast.ToInt(coalesce("REDIS_DB", 0))

	config.KafkaBrokers = cast.ToStringSlice(coalesce("KAFKA_BROKERS", []string{"kafka:9092"}))

	config.LOG_PATH = cast.ToString(coalesce("LOG_PATH", "logs/info.log"))

	return config
}

func coalesce(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
