package transcoder

// Multi Rendition HLS transcoding

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/harshjoeyit/video-streaming/storage"
)

// TranscodeToMultiRenditionHLS transcodes MP4 video file to HLS format
// with multiple renditions (e.g., 720p, 480p, 240p)
func TranscodeToMultiRenditionHLS(fileId string) error {
	src := storage.GetUploadedVideoPath(fileId)
	dstDir := storage.GetChunkedVideoPath(fileId)

	if err := storage.CreateDirectoryIfNotExists(dstDir); err != nil {
		return fmt.Errorf("failed to create dir for segments: %w", err)
	}

	masterFile := "master.m3u8"
	progPath := fmt.Sprintf("%s/v%%v/prog.m3u8", dstDir)
	segPatternPath := fmt.Sprintf("%s/v%%v/seg_%%03d.ts", dstDir)

	// Note: ffmpeg-go library is not consistent with multiple same
	// named flags, for eg. -map. Hence using cmd.exec()
	cmd := exec.CommandContext(context.Background(), "ffmpeg",
		"-i", src,
		// video+audio mappings
		"-map", "0:v", "-map", "0:a",
		"-map", "0:v", "-map", "0:a",
		"-map", "0:v", "-map", "0:a",
		// video encoding settings (for 3 streams)
		"-c:v", "libx264",
		"-s:v:0", "426x240", "-b:v:0", "400k",
		"-s:v:1", "854x480", "-b:v:1", "800k",
		"-s:v:2", "1280x720", "-b:v:2", "2000k",
		// audio encoding settings (for 3 streams)
		"-c:a", "aac",
		"-b:a:0", "96k",
		"-b:a:1", "96k",
		"-b:a:2", "96k",
		// HLS output options
		"-f", "hls",
		"-hls_time", "6",
		// Video on Demand
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", segPatternPath,
		"-master_pl_name", masterFile,
		// // map the streams to names (240p, 480p, 720p)
		"-var_stream_map", "v:0,a:0,name:240p v:1,a:1,name:480p v:2,a:2,name:720p",
		progPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
