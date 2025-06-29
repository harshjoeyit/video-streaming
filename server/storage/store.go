package storage

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Store the uploaded video file. Read and return manifest and segments

const (
	uploadedVideosPath = "uploaded_videos"
	chunkedVideosPath  = "chunked_videos"
)

func InitVideoStorage() error {
	// Create directories if they don't exist
	if err := CreateDirectoryIfNotExists(uploadedVideosPath); err != nil {
		return err
	}
	if err := CreateDirectoryIfNotExists(chunkedVideosPath); err != nil {
		return err
	}

	return nil
}

// UploadVideo saves the uploaded video file to the server
func UploadVideo(
	file *multipart.FileHeader,
	fileID string,
	c *gin.Context,
) error {
	filePath := filepath.Join(uploadedVideosPath, fmt.Sprintf("%s.mp4", fileID))

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// GetProcessedVideoList returns a list of processed video directories
func GetProcessedVideoList() ([]string, error) {
	files, err := os.ReadDir(chunkedVideosPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read videos: %w", err)
	}

	fmt.Printf("Found %d processed video directories, %v\n", len(files), files)

	videoList := []string{}
	for _, file := range files {
		if file.IsDir() {
			videoList = append(videoList, file.Name())
		}
	}

	return videoList, nil
}

// GetUploadedVideoPath returns the path where the uploaded video is stored
func GetUploadedVideoPath(fileID string) string {
	return filepath.Join(uploadedVideosPath, fmt.Sprintf("%s.mp4", fileID))
}

// GetChunkedVideoPath returns the path where the manifest and segments are stored
func GetChunkedVideoPath(fileID string) string {
	return filepath.Join(chunkedVideosPath, fileID)
}

// GetVideoManifestPath returns the path for the video manifest file (m3u8)
// This is used for transcoding to HLS format
func GetVideoManifestPath(fileID string) (string, error) {
	path := filepath.Join(GetChunkedVideoPath(fileID), "playlist.m3u8")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("manifest file does not exist: %s", path)
	} else if err != nil {
		return "", fmt.Errorf("error checking manifest file: %w", err)
	}

	return path, nil
}

// GetVideoSegmentPath returns the path for a specific video segment
func GetVideoSegmentPath(fileID, segment string) (string, error) {
	path := filepath.Join(GetChunkedVideoPath(fileID), segment)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("segment file does not exist: %s", path)
	} else if err != nil {
		return "", fmt.Errorf("error checking segment file: %w", err)
	}

	return path, nil
}

// GetVideoSegmentPatternPath returns the pattern path for video segments
// This is used for transcoding to HLS format
func GetVideoSegmentPatternPath(fileID string) string {
	return filepath.Join(GetChunkedVideoPath(fileID), "segment_%03d.ts")
}
func CreateDirectoryIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}
