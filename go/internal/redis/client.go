package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	// Allow host and port to be overridden by environment variables
	host := getEnv("REDIS_HOST", "redis")
	port := getEnv("REDIS_PORT", "6379")
	addr := fmt.Sprintf("%s:%s", host, port)

	// Create new Redis client
	RDB = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// Test connection
	if err := RDB.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", addr, err)
	}

	log.Println("✅ Connected to Redis at", addr)
}

// getEnv returns the value of the env variable or fallback if not set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
