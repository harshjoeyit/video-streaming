package transcoder

// Multi Rendition HLS transcoding

import (
	"fmt"

	"github.com/harshjoeyit/video-streaming/storage"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// TranscodeToMultiRenditionHLS transcodes MP4 video file to HLS format
// with multiple renditions (e.g., 720p, 480p, 240p)
func TranscodeToMultiRenditionHLS(fileId string) error {
	src := storage.GetUploadedVideoPath(fileId)
	dstDir := storage.GetChunkedVideoPath(fileId)

	if err := storage.CreateDirectoryIfNotExists(dstDir); err != nil {
		return fmt.Errorf("failed to create dir for segments: %w", err)
	}

	manifestFile, err := storage.GetVideoManifestPath(fileId)
	if err != nil {
		return fmt.Errorf("failed to get manifest file path: %w", err)
	}

	segPatternPath := storage.GetVideoSegmentPatternPath(fileId)

	return ffmpeg.
		Input(src).
		Output(
			manifestFile,
			kwargsWithFixedBitrate(segPatternPath),
		).
		OverWriteOutput().
		GlobalArgs("-progress", "pipe:2", "-nostats").
		ErrorToStdOut().
		Run()
}
