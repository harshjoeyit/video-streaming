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

	segPatternPath := fmt.Sprintf("%s/v%%v/seg_%%03d.ts", dstDir)

	return ffmpeg.
		Input(src).
		Output(
			fmt.Sprintf("%s/v%%v/prog.m3u8", dstDir),
			ffmpeg.KwArgs{
				"filter_complex": "[0:v]split=3[v1][v2][v3];" +
					"[v1]scale=w=426:h=240[v1out];" +
					"[v2]scale=w=854:h=480[v2out];" +
					"[v3]scale=w=1280:h=720[v3out]",
				"map":                  "[v1out]",
				"map:1":                "[v2out]",
				"map:2":                "[v3out]",
				"map:a":                "0:a",
				"c:v:0":                "libx264",
				"b:v:0":                "400k",
				"c:v:1":                "libx264",
				"b:v:1":                "800k",
				"c:v:2":                "libx264",
				"b:v:2":                "2000k",
				"c:a":                  "aac",
				"b:a":                  "96k",
				"f":                    "hls",
				"hls_time":             "6",
				"hls_playlist_type":    "vod",
				"master_pl_name":       "master.m3u8",
				"hls_segment_filename": segPatternPath,
				"var_stream_map":       "v:0,a:0 v:1,a:0 v:2,a:0",
			},
		).
		OverWriteOutput().
		GlobalArgs("-progress", "pipe:2", "-nostats").
		ErrorToStdOut().
		Run()
}
