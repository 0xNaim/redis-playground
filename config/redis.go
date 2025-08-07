package config

import (
	"context"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// Initialize and returns a Redis client
func InitRedis() *redis.Client {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic("No .env file found, using default configuration")
	}

	// Get Redis configuration from environment variables or use defaults
	addr := getEnv("REDIS_ADDR", "localhost:6379")
	password := getEnv("REDIS_PASSWORD", "")
	db := getEnvAsInt("REDIS_DB", 0)

	// Create a new Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Create a context for the Redis operations
	ctx := context.Background()
	if ctx == nil {
		panic("Failed to create Redis context")
	}

	// Test the connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	// Log the successful connection
	println("Connected to Redis at", addr)
	return client
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
