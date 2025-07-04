package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/harshjoeyit/video-streaming/storage"
	"github.com/harshjoeyit/video-streaming/transcoder"
)

// Todo: pre‑signed URL
func uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	if file.Header.Get("Content-Type") != "video/mp4" &&
		file.Header.Get("Content-Type") != "multipart/form-data" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Only MP4 files are allowed"},
		)
		return
	}

	// Single rendition or multi-rendition (ABR)
	rendition := c.PostForm("rendition")
	if rendition != "single" && rendition != "multi" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid rendition type. Use 'single' or 'multi'."},
		)
		return
	}

	// Generate a unique file ID using UUID
	fileID := uuid.New().String()

	if err := storage.UploadVideo(file, fileID, c); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	// Set status to "PROCESSING" in Redis
	if err := storage.SetVideoProcStatus(fileID, storage.StatusProcessing); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to set file status in Redis"},
		)
		return
	}

	// Start transcoding the video to HLS format in a separate goroutine
	// Allows the server to respond immediately while processing continues
	go func(fileID string) {
		switch rendition {
		case transcoder.SingleRendition:
			if err := transcoder.TranscodeToHLS(fileID); err != nil {
				log.Printf("failed to transcode video: %v", err)

				// Set status to "FAILED" in Redis
				if err := storage.SetVideoProcStatus(fileID, storage.StatusFailed); err != nil {
					log.Printf("failed to set file status in Redis: %v", err)
				}

				return
			}
		case transcoder.MultiRendition:
			// if err := transcoder.TranscodeToHLS(fileID); err != nil {
			if err := transcoder.TranscodeToMultiRenditionHLS(fileID); err != nil {
				log.Printf("failed to transcode video: %v", err)

				// Set status to "FAILED" in Redis
				if err := storage.SetVideoProcStatus(fileID, storage.StatusFailed); err != nil {
					log.Printf("failed to set file status in Redis: %v", err)
				}

				return
			}
		}

		// Set status to "READY" in Redis
		if err := storage.SetVideoProcStatus(fileID, storage.StatusReady); err != nil {
			log.Printf("failed to set file status in Redis: %v", err)
		}
	}(fileID)

	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  "success",
			"id":      fileID,
			"message": "File uploaded successfully. Processing started.",
		},
	)
}

func statusHandler(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	status, err := storage.GetVideoStatus(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": fileID, "status": status})
}

func assetsHandler(c *gin.Context) {
	videoList, err := storage.GetProcessedVideoList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"videos": videoList})
}

func playlistHandler(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	playlistPath, err := storage.GetVideoManifestPath(videoID)
	if err != nil {
		log.Printf("failed to serve video manifest: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Serve the manifest file
	c.File(playlistPath)
}

func segmentHandler(c *gin.Context) {
	videoID := c.Param("id")
	segment := c.Param("segment")
	if videoID == "" || segment == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Video ID and segment are required"},
		)
		return
	}

	segmentPath, err := storage.GetVideoSegmentPath(videoID, segment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Serve the segment file
	c.File(segmentPath)
}

// Adaptive Bitrate Streaming (ABR) handlers

// playlistHandlerABR serves master file
func playlistHandlerABR(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
		return
	}

	masterPlaylistFile, err := storage.GetVideoManifestPathABR(videoID)
	if err != nil {
		log.Printf("failed to serve master video manifest: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Serve the manifest file
	c.File(masterPlaylistFile)
}

func renditionPlaylistHandler(c *gin.Context) {
	videoID := c.Param("id")
	rendition := c.Param("rendition")
	if videoID == "" || rendition == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Video ID and rendition are required"},
		)
		return
	}

	playlistPath, err := storage.GetRenditionPlaylistPath(videoID, rendition)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	// Serve the rendition playlist file (eg. /v0/prog.m3u8)
	c.File(playlistPath)
}

func segmentHandlerABR(c *gin.Context) {
	videoID := c.Param("id")
	rendition := c.Param("rendition")
	segment := c.Param("segment")
	if videoID == "" || rendition == "" || segment == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Video ID, rendition and segment are required"},
		)
		return
	}

	segmentPath, err := storage.GetVideoSegmentPathABR(videoID, rendition, segment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Serve the segment file
	c.File(segmentPath)
}
