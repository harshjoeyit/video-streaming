package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	uploadedVideosPath = "uploaded_videos"
	chunkedVideosPath  = "chunked_videos"
)

func main() {
	ge := gin.Default()

	ge.Use(CORSMiddleware())

	// Create directories if they don't exist
	if err := os.MkdirAll(uploadedVideosPath, 0755); err != nil {
		log.Fatalf("failed to create uploaded_videos directory: %v", err)
	}
	if err := os.MkdirAll(chunkedVideosPath, 0755); err != nil {
		log.Fatalf("failed to create chunked_videos directory: %v", err)
	}

	// Upload a MP4 or multipart video file
	// Todo: pre‑signed URL
	ge.POST("/upload", func(c *gin.Context) {
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

		// generate a unique folder name (UUID)
		fileID := uuid.New().String()
		filePath := filepath.Join(uploadedVideosPath, fmt.Sprintf("%s.mp4", fileID))

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Failed to save file"},
			)
			return
		}

		err = TranscodeToHLS(fileID)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": fmt.Sprintf("Failed to transcode video: %v", err)},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{
				"status": "success",
				"id":     fileID,
			},
		)
	})

	// Return list of all available videos
	ge.GET("/assets", func(c *gin.Context) {
		files, err := os.ReadDir(chunkedVideosPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read videos"})
			return
		}

		videoList := []string{}
		for _, file := range files {
			if file.IsDir() {
				videoList = append(videoList, file.Name())
			}
		}

		c.JSON(http.StatusOK, gin.H{"videos": videoList})
	})

	// Serve the HLS manifest file (playlist.m3u8) for a video
	ge.GET("/assets/:id/playlist.m3u8", func(c *gin.Context) {
		videoID := c.Param("id")
		if videoID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Video ID is required"})
			return
		}

		playlistPath := filepath.Join(chunkedVideosPath, videoID, "playlist.m3u8")
		if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
			return
		}

		c.File(playlistPath)
	})

	// Serve individual video segments
	ge.GET("/assets/:id/:segement", func(c *gin.Context) {
		videoID := c.Param("id")
		segment := c.Param("segement")
		if videoID == "" || segment == "" {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "Video ID and segment are required"},
			)
			return
		}

		segmentPath := filepath.Join(chunkedVideosPath, videoID, segment)
		if _, err := os.Stat(segmentPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Segment not found"})
			return
		}

		// Serve the segment file directly
		c.File(segmentPath)
	})

	ge.Run(":9090")
}

// TranscodeToHLS transcodes the MP4 video file to HLS format.
// This creates the manifest file (m3u8) and the segment files (ts).
// Ref: https://ffmpeg.org/ffmpeg-formats.html#hls-2
//
// Note: This should ideally be run via a queue on separate worker nodes
func TranscodeToHLS(fileId string) error {
	// Source file path
	src := filepath.Join(uploadedVideosPath, fmt.Sprintf("%s.mp4", fileId))
	// Destination directory for manifest and segments
	dstDir := filepath.Join(chunkedVideosPath, fileId)

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create dir for segments: %w", err)
	}

	manifestFile := filepath.Join(dstDir, "playlist.m3u8")
	segPattern := filepath.Join(dstDir, "segment_%03d.ts")

	return ffmpeg.
		Input(src).
		Output(
			manifestFile,
			kwargsWithCRFCompression(segPattern),
		).
		OverWriteOutput().
		GlobalArgs("-progress", "pipe:2", "-nostats").
		ErrorToStdOut().
		Run()
}

// TranscodeToMultiRenditionHLS transcodes MP4 video file to HLS format
// with multiple renditions (e.g., 720p, 480p, 240p)
func TranscodeToMultiRenditionHLS(fileId string) error {
	// Source file path
	src := filepath.Join(uploadedVideosPath, fmt.Sprintf("%s.mp4", fileId))
	// Destination directory for manifest and segments
	dstDir := filepath.Join(chunkedVideosPath, fileId)

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create dir for segments: %w", err)
	}

	manifestFile := filepath.Join(dstDir, "playlist.m3u8")
	segPattern := filepath.Join(dstDir, "segment_%03d.ts")

	return ffmpeg.
		Input(src).
		Output(
			manifestFile,
			kwargsWithFixedBitrate(segPattern),
		).
		OverWriteOutput().
		GlobalArgs("-progress", "pipe:2", "-nostats").
		ErrorToStdOut().
		Run()
}

// kwargsWithFixedBitrate returns ffmpeg.KwArgs for HLS transcoding
// with fixed bitrate encoding for video of 5 Mbps.
func kwargsWithFixedBitrate(segPattern string) ffmpeg.KwArgs {
	k := getDefaultHLSKwargs(segPattern)
	k["b:v"] = "5M"      // target video bitrate.
	k["maxrate"] = "5M"  // maximum video bitrate.
	k["bufsize"] = "10M" // rate control buffer size.

	return k
}

// kwargsWithCRFCompression returns ffmpeg.KwArgs for HLS transcoding
// with CRF compression of 26.
func kwargsWithCRFCompression(segPattern string) ffmpeg.KwArgs {
	k := getDefaultHLSKwargs(segPattern)
	// crf 26: compress more, but quality is worse (default = crf 23)
	k["crf"] = "26"
	// preset slow: slower encode, better compression (safe for offline VoD)
	k["preset"] = "slow"

	return k
}

func getDefaultHLSKwargs(segPattern string) ffmpeg.KwArgs {
	return ffmpeg.KwArgs{
		// Video codec: re-encodes stream with x264 H.264 encoder
		"c:v": "libx264",
		// Main profile for modern devices and web streaming
		"profile:v": "main",
		// Caps encoder complexity at Level 3.1 (≈ 720p @ 30 fps, ≤14 Mb/s)
		"level": "3.1",
		// audio codec - AAC is required by HLS
		"c:a": "aac",
		// audio bitrate - 128 kb/s which is good enough for VoD
		"b:a": "128k",
		// format is HLS - creates manifect file and segments
		"f": "hls",
		// segment duration in seconds
		"hls_time": "6",
		// video on demand. Other options are "live" and "event"
		"hls_playlist_type":    "vod",
		"hls_segment_filename": segPattern,
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Todos:
// 1. Multiple renditions of the same video (e.g., 720p, 480p, 360p)
// 2. Using S3 and CDN for video storage and delivery
