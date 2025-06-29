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
	uploadedVideosPath = "assets/uploaded_videos"
	chunkedVideosPath  = "assets/chunked_videos"
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

type VideoItem struct {
	ID  string `json:"id"`
	ABR bool   `json:"abr"`
}

// GetProcessedVideoList returns a list of processed video directories
func GetProcessedVideoList() ([]VideoItem, error) {
	files, err := os.ReadDir(chunkedVideosPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read videos: %w", err)
	}

	videoList := []VideoItem{}
	for _, file := range files {
		if file.IsDir() {
			item := VideoItem{
				ID:  file.Name(),
				ABR: false,
			}

			masterPath := filepath.Join(chunkedVideosPath, file.Name(), "master.m3u8")
			if _, err := os.Stat(masterPath); err == nil {
				item.ABR = true
			}

			videoList = append(videoList, item)
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

// getVideoManifestPath returns the path for the video manifest file (m3u8)
// based on the rendition type.
func getVideoManifestPath(fileID string, rendition string) string {
	// For single rendition HLS, the manifest file is named "playlist.m3u8"
	switch rendition {
	case "single":
		return filepath.Join(GetChunkedVideoPath(fileID), "playlist.m3u8")
	case "multi":
		return filepath.Join(GetChunkedVideoPath(fileID), "master.m3u8")
	default:
		return ""
	}
}

// ServeVideoManifest returns the path for the video manifest file
// if exists, otherwise return error
func GetVideoManifestPath(fileID string) (string, error) {
	path := getVideoManifestPath(fileID, "single")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("manifest file does not exist for file ID %s: %w", fileID, err)
	} else if err != nil {
		return "", fmt.Errorf("failed to get manifest for file ID %s: %w", fileID, err)
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

// For Adaptive Bitrate Streaming

func GetVideoManifestPathABR(fileID string) (string, error) {
	path := getVideoManifestPath(fileID, "multi")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("master file does not exist for file ID %s: %w", fileID, err)
	} else if err != nil {
		return "", fmt.Errorf("failed to get master for file ID %s: %w", fileID, err)
	}

	return path, nil
}

func GetRenditionPlaylistPath(fileID, rendition string) (string, error) {
	path := filepath.Join(GetChunkedVideoPath(fileID), rendition, "prog.m3u8")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("rendition playlist file does not exist: %s", path)
	} else if err != nil {
		return "", fmt.Errorf("error checking rendition playlist file: %w", err)
	}

	return path, nil
}

// GetVideoSegmentPathABR returns the path for a specific video segment
// for the given file and rendition
func GetVideoSegmentPathABR(fileID, rendition, segment string) (string, error) {
	path := filepath.Join(GetChunkedVideoPath(fileID), rendition, segment)
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
