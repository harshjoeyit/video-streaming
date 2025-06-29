package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	StatusProcessing = "PROCESSING"
	StatusReady      = "READY"
	StatusFailed     = "FAILED"
)

// SetVideoProcStatus sets the processing status of a video in Redis
func SetVideoProcStatus(fileID, status string) error {
	if rdb == nil {
		return fmt.Errorf("redis client not initialized")
	}

	key := getStatusKey(fileID)
	if err := rdb.Set(context.Background(), key, status, 0).Err(); err != nil {
		return fmt.Errorf("failed to set video status: %w", err)
	}

	return nil
}

// GetVideoStatus retrieves the processing status of a video from Redis
func GetVideoStatus(fileID string) (string, error) {
	if rdb == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	key := getStatusKey(fileID)
	status, err := rdb.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("video status not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to get video status: %w", err)
	}

	return status, nil
}

func getStatusKey(fileID string) string {
	return fmt.Sprintf("video:%s:status", fileID)
}
