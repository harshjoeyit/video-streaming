package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(ge *gin.Engine) {
	// Upload of a MP4 or multipart video file
	ge.POST("/upload", uploadHandler)

	// Status of the uploaded video file (e.g., PROCESSING, READY, FAILED)
	ge.GET("/status/:id", statusHandler)

	// Return list of all available videos
	ge.GET("/assets", assetsHandler)

	// Serve the HLS manifest file (playlist.m3u8) for a video
	ge.GET("/assets/:id/playlist.m3u8", playlistHandler)

	// Serve individual video segments
	ge.GET("/assets/:id/:segment", segmentHandler)
}

func UseCORS(ge *gin.Engine) {
	// CORS middleware to allow cross-origin requests
	ge.Use(corsMiddleware())
}
