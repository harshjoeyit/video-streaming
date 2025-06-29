Video streaming prototype for Video on demand

Requirements:
go: 1.24.1
ffmpeg binary installed and set in $PATH
Used: 7.7.1

1. Packaging protocol: Http Live Stream (HLS)
2. Video codec: H.264 + AAC
3. Segment length: Six‑second segments (-hls_time 6) are a safe default.

Basics:
- https://howvideo.works/ 

- What is m3u8 file?
    - Playlist: M3U files are playlists, meaning they organize and list media files. 
    - Text-based: They are plain text files, typically with a .m3u or .m3u8 extension. 
    - References, not media: M3U files contain paths to media files (audio or video) on your computer or on the internet, not the media data itself. 
    - Common uses: M3U files are used for streaming internet radio, creating custom playlists for music or videos, and in HTTP Live Streaming (HLS)


Bitrate: Data transferred per unit time.
File size = bitrate * duration
CBR vs VBR: most live streaming uses CBR
Codecs: Coder / Decoder - kind of like a language
    Video codec: 
        H.246 is a compression standard and another name for Advanced Video Coding (AVC). It's a lossy compression method that reduces file sizes while maintaining good video quality, making it ideal for various applications like streaming, Blu-ray discs, and broadcast television. 
        H.265 (HEVC):
        Pros: Up to 50% better compression efficiency than H.264, allowing for smaller file sizes or better quality at the same size, ideal for streaming high-definition content. 
        Cons: Requires more processing power, potentially impacting performance on older devices, compatibility is not as widespread as H.264
    AAC - Advanced Audio Coding, is a digital audio compression format known for its efficiency and higher quality compared to MP3 at the same bitrates.
    Bitrate is used during encoding which affects the quality of the footage. Higher bitrate would result in better quality but more file size.
    Container (MOV, MP4, MKV) contains the Codec file.

What is MPEG-TS?
    MPEG-TS, or MPEG Transport Stream, is a digital container which multiplexes multiple audio and video streams into a single bitstream. 
    - It's known for its ability to maintain stream integrity even in unreliable environments like the internet or satellite links.
    - It uses small packets, each with a header containing information like timing and sequence, to ensure proper reassembly at the receiving end. 
    - MPEG-TS includes features for error correction and synchronization, which is crucial for maintaining a stable stream, especially in environments prone to packet loss or network congestion
    - It is optimized (and widely used) for live streaming, digital TV broadcasts

MPEG - Moving Pictures Expert Group

https://ffmpeg.org/ffmpeg-formats.html#hls-2


For compression along with segmenetation:

- Packaging ≠ compression. If you want smaller files, dial the encoder, not the segmenter.
- Expect only a ~5 % bump when moving from MP4 to HLS‑TS at normal streaming bit‑rates.

1. The CRF (Constant Rate Factor) mode is ideal if:
    You want good visual quality
    But you don’t care about exact bitrates
    And you’re okay with slight variation in size per video
    Use case: Video on Demand (VoD)

    ```
    -c:v libx264 -crf 26 -preset slow
    ```

    CRF ranges:
    18–20 → visually lossless (high quality, large size)
    23 (default) → decent quality, medium size
    28–32 → much smaller file, visibly compressed

2. Set Explicit Bitrate (if you want consistent size)
    Instead of CRF, you can target a specific bitrate, for example:

    ```
    -c:v libx264 -b:v 1M -maxrate 1M -bufsize 2M
    ```
    1M = 1 megabit/sec (≈ 450 MB/hour video)
    This gives you more predictable file sizes, especially useful if you're billing per GB on CDN.
    Use case: Live streaming


 3. Downscale Resolution (big size savings)
    If you're streaming only to mobile/web, reducing resolution gives massive savings:
    854x480	Saves ~60% space vs 1080p

    1M = 1 megabit/sec (≈ 450 MB/hour video)

4. Use a More Efficient Codec (optional for modern clients)
    H.265 (HEVC) compresses better than H.264, ~30–50% smaller at same quality. Con: limited support


What is a `.m3u8` file?

- `.m3u8` files are UTF-8 encoded playlists used in HLS to stream media by dividing it into smaller `.ts` (MPEG-TS) segments.
- These files list the metadata and segments required for playback.

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

More details - https://chatgpt.com/share/6860c089-b284-8012-8a27-b456c24584df

More details - https://chatgpt.com/share/6860c0a0-0e94-8012-a930-aaadaeed404b
