package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/harshjoeyit/video-streaming/handler"
	"github.com/harshjoeyit/video-streaming/storage"
)

func main() {
	// Initialize Redis client
	storage.InitRedisClient()
	defer storage.CloseRedisClient()

	// Initialize video storage directories
	if err := storage.InitVideoStorage(); err != nil {
		log.Fatalf("failed to initialize video storage: %v", err)
	}

	ge := gin.Default()

	handler.UseCORS(ge)
	handler.RegisterRoutes(ge)

	ge.Run(":9090")
}

// Todos:
// 1. Multiple renditions of the same video (e.g., 720p, 480p, 360p)
// 2. Using S3 and CDN for video storage and delivery
