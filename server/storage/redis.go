package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	addr     = "localhost:6379"
	password = "" // no password set
)

// Redis client connection
var rdb *redis.Client

func InitRedisClient() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	// Test the connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// GetRedisClient returns the initialized Redis client
func GetRedisClient() *redis.Client {
	if rdb == nil {
		panic("Redis client not initialized. Call InitRedisClient first.")
	}
	return rdb
}

// CloseRedisClient closes the Redis client connection
func CloseRedisClient() error {
	if rdb == nil {
		return nil // already closed or not initialized
	}

	if err := rdb.Close(); err != nil {
		return fmt.Errorf("failed to close Redis client: %w", err)
	}
	rdb = nil // reset the client
	return nil
}
