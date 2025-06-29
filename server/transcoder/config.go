package transcoder

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

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
