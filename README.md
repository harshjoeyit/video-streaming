# Video Streaming Prototype for Video on Demand (VoD)

## Requirements
- Go: 1.24.1
- ffmpeg (v7.7.1) binary installed and set in `$PATH`

---

## Features

1. **Packaging protocol:** HTTP Live Streaming (HLS)
2. **Video codec:** H.264 + AAC
3. **Segment length:** Six‑second segments (`-hls_time 6`) are a safe default.

---

## API Docs

[Postman Collection Link](https://web.postman.co/workspace/My-Workspace~52c68d62-65f0-4416-8915-0c154d43c09b/collection/13052768-43235c2a-53e0-458b-970b-8fa38f9a16bb?action=share&source=copy-link&creator=13052768)

---

## Basics

### [How Video Works](https://howvideo.works/)

### What is an m3u8 file?

- **Playlist:** M3U files are playlists, meaning they organize and list media files.
- **Text-based:** They are plain text files, typically with a `.m3u` or `.m3u8` extension.
- **References, not media:** M3U files contain paths to media files (audio or video) on your computer or on the internet, not the media data itself.
- **Common uses:** M3U files are used for streaming internet radio, creating custom playlists for music or videos, and in HTTP Live Streaming (HLS).

---

### Video Concepts

- **Bitrate:** Data transferred per unit time.  
  `File size = bitrate * duration`
- **CBR vs VBR:** Most live streaming uses CBR (Constant Bitrate).
- **Codecs:** Coder / Decoder - kind of like a language.
    - **Video codec:**  
      - **H.264:** Also known as AVC, a lossy compression standard for good quality and small file sizes.
      - **H.265 (HEVC):** Up to 50% better compression than H.264, but requires more processing power and has less compatibility.
    - **AAC:** Advanced Audio Coding, more efficient and higher quality than MP3 at the same bitrates.
- **Container:** (MOV, MP4, MKV) contains the codec file.

---

### What is MPEG-TS?

- MPEG-TS (MPEG Transport Stream) is a digital container that multiplexes multiple audio and video streams into a single bitstream.
- Maintains stream integrity in unreliable environments (internet, satellite).
- Uses small packets with headers for timing and sequence.
- Includes error correction and synchronization.
- Widely used for live streaming and digital TV broadcasts.

---

### Compression & Segmentation

- **Packaging ≠ compression:** If you want smaller files, adjust the encoder, not the segmenter.
- **Expect only a ~5% bump** when moving from MP4 to HLS‑TS at normal streaming bit‑rates.

---

### Encoding Tips

#### 1. **CRF (Constant Rate Factor) Mode**

- Ideal if you want good visual quality and don't care about exact bitrates.
- Use case: Video on Demand (VoD)
- Example:
    ```sh
    -c:v libx264 -crf 26 -preset slow
    ```
- **CRF ranges:**
    - 18–20: Visually lossless (high quality, large size)
    - 23 (default): Decent quality, medium size
    - 28–32: Much smaller file, visibly compressed

#### 2. **Set Explicit Bitrate**

- For consistent file size, use:
    ```sh
    -c:v libx264 -b:v 1M -maxrate 1M -bufsize 2M
    ```
- 1M = 1 megabit/sec (≈ 450 MB/hour video)
- Use case: Live streaming

#### 3. **Downscale Resolution**

- Reducing resolution gives massive savings.
- Example: 854x480 saves ~60% space vs 1080p.

#### 4. **Use a More Efficient Codec**

- H.265 (HEVC) compresses better than H.264 (~30–50% smaller at same quality).
- Con: limited support.

---

### What is a `.m3u8` file?

- `.m3u8` files are UTF-8 encoded playlists used in HLS to stream media by dividing it into smaller `.ts` (MPEG-TS) segments.
- These files list the metadata and segments required for playback.

**Example (Single Rendition):**

```
#EXTM3U                     // Format .m3u (required for HLS)
#EXT-X-VERSION:3            // Version
#EXT-X-TARGETDURATION:8     // Maximum segment duration (in secs). Useful for players to buffer
#EXT-X-MEDIA-SEQUENCE:0     // Sequence number of 1st segement 
#EXT-X-PLAYLIST-TYPE:VOD    // Video on demand (denotes file is complete and won't change)
#EXTINF:6.350978,           // #EXTINF:<duration> of segment (in secs)
segment_000.ts              // segment filename
#EXTINF:5.916433,
segment_001.ts
#EXTINF:8.089133,
segment_002.ts
#EXT-X-ENDLIST              // End of playlist

```

**Example (Multipel Renditions - ABR):**

```
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:BANDWIDTH=602081,AVERAGE-BANDWIDTH=573792,RESOLUTION=426x240,CODECS="avc1.64001e,mp4a.40.2"
v240p/prog.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=1062305,AVERAGE-BANDWIDTH=978475,RESOLUTION=854x480,CODECS="avc1.64001f,mp4a.40.2"
v480p/prog.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=2485931,AVERAGE-BANDWIDTH=2190930,RESOLUTION=1280x720,CODECS="avc1.640020,mp4a.40.2"
v720p/prog.m3u8
```
***v720p/prog.m3u8***
```
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:8
#EXT-X-MEDIA-SEQUENCE:0
#EXT-X-PLAYLIST-TYPE:VOD
#EXTINF:8.333333,
seg_000.ts
#EXTINF:4.166667,
seg_001.ts
#EXTINF:8.333333,
seg_002.ts
#EXT-X-ENDLIST
```


### References

- [ffmpeg HLS documentation](https://ffmpeg.org/ffmpeg-formats.html#hls-2)
- [ChatGPT Explanation 1](https://chatgpt.com/share/6860c089-b284-8012-8a27-b456c24584df)
- [ChatGPT Explanation 2](https://chatgpt.com/share/6860c0a0-0e94-8012-a930-aaadaeed404b)