package transcoder

import (
	"fmt"
	"path/filepath"

	"github.com/harshjoeyit/video-streaming/storage"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	SingleRendition = "single"
	MultiRendition  = "multi"
)

// Single Rendition HLS transcoding

// TranscodeToHLS transcodes the MP4 video file to HLS format.
// This creates the manifest file (m3u8) and the segment files (ts).
// Ref: https://ffmpeg.org/ffmpeg-formats.html#hls-2
//
// Note: This should ideally be run via a queue on separate worker nodes
func TranscodeToHLS(fileId string) error {
	src := storage.GetUploadedVideoPath(fileId)
	dstDir := storage.GetChunkedVideoPath(fileId)

	if err := storage.CreateDirectoryIfNotExists(dstDir); err != nil {
		return fmt.Errorf("failed to create dir for segments: %w", err)
	}

	manifestFile := filepath.Join(dstDir, "playlist.m3u8")
	segPatternPath := storage.GetVideoSegmentPatternPath(fileId)

	return ffmpeg.
		Input(src).
		Output(
			manifestFile,
			kwargsWithCRFCompression(segPatternPath),
		).
		OverWriteOutput().
		GlobalArgs("-progress", "pipe:2", "-nostats").
		ErrorToStdOut().
		Run()
}
