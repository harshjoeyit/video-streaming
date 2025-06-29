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

	// Below APIs are for single rendition playback (no-ABR)

	// No-ABR: Serve the HLS manifest file (playlist.m3u8) for a video
	ge.GET("/assets/:id/playlist.m3u8", playlistHandler)

	// No-ABR: Serve individual video segments
	ge.GET("/assets/:id/:segment", segmentHandler)

	// Below APIs are for multi-rendition playback (ABR)

	// ABR: Serve the master file (master.m3u8) for a video
	ge.GET("/assets/abr/:id/master.m3u8", playlistHandlerABR)

	// ABR: Serve a rendition manifest (for Adaptive Bitrate Streaming)
	ge.GET("/assets/abr/:id/:rendition/prog.m3u8", renditionPlaylistHandler)

	// ABR: Serve a rendition manifest (for Adaptive Bitrate Streaming)
	ge.GET("/assets/abr/:id/:rendition/:segment", segmentHandlerABR)

}

func UseCORS(ge *gin.Engine) {
	// CORS middleware to allow cross-origin requests
	ge.Use(corsMiddleware())
}
