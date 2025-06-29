import { useRef, useEffect, useState } from "react";
import "./App.css";
import Hls from "hls.js";

const baseUrl = "http://localhost:9090";

export default function App() {
  // Call /assets API to list of available videos
  const [videos, setVideos] = useState([]);
  const [selectedVideo, setSelectedVideo] = useState(null);

  useEffect(() => {
    fetch(`${baseUrl}/assets`)
      .then((response) => response.json())
      .then((data) => {
        console.log("videos", data);
        setVideos(data.videos || []);
      })
      .catch((error) => console.error("Error fetching videos:", error));
  }, []);

  return (
    <div className="App">
      <h2>Video Streaming App</h2>
      <div
        style={{ display: "flex", flexDirection: "row", alignItems: "center" }}
      >
        <VideoList videos={videos} onSelect={setSelectedVideo} />
        {selectedVideo && <HLSVideoPlayer src={getVideoSrc(selectedVideo)} />}
      </div>
    </div>
  );
}

function VideoList({ videos, onSelect }) {
  return (
    <div>
      <h4>Available Videos</h4>
      <ol>
        {videos.map((video) => (
          <li key={video.id} style={{ fontSize: "14px" }}>
            <button onClick={() => onSelect(video)}>
              {video.id}
              {video.abr ? " | ABR" : " | No-ABR"}
            </button>
          </li>
        ))}
      </ol>
    </div>
  );
}

function HLSVideoPlayer({ src }) {
  const videoRef = useRef(null);
  const [levels, setLevels] = useState([]); // [{id, label}]
  const [hlsObj, setHlsObj] = useState(null);

  useEffect(() => {
    if (canPlayNativeHls) {
      // ⇢ Safari
      const v = videoRef.current;
      if (v) v.src = src;
      return;
    }

    if (Hls.isSupported() && videoRef.current) {
      // ⇢ hls.js path
      const hls = new Hls();
      hls.loadSource(src);
      hls.attachMedia(videoRef.current);

      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        const opts = hls.levels.map((l, i) => ({
          id: i,
          label: l.height
            ? `${l.height}p`
            : `${Math.round(l.bitrate / 1000)} kbps`,
        }));
        setLevels([{ id: -1, label: "Auto" }, ...opts]);
      });

      setHlsObj(hls);
      return () => hls.destroy();
    }
  }, [src]);

  // Change resil
  const handleChangeResolution = (e) => {
    const level = parseInt(e.target.value, 10);
    if (hlsObj) hlsObj.currentLevel = level; // -1 = Auto
  };

  return (
    <div>
      <h4>Playing: {src}</h4>
      <video ref={videoRef} controls style={{ width: "400px" }} />
      {/* Resolutions */}
      {!canPlayNativeHls && levels.length > 1 && (
        <select onChange={handleChangeResolution} style={{ marginTop: 8 }}>
          {levels.map((l) => (
            <option key={l.id} value={l.id}>
              {l.label}
            </option>
          ))}
        </select>
      )}
    </div>
  );
}

const canPlayNativeHls = (() => {
  const video = document.createElement("video");
  return video.canPlayType("application/vnd.apple.mpegurl");
})();

function getVideoSrc(video) {
  if (video.abr) {
    return `${baseUrl}/assets/abr/${video.id}/master.m3u8`;
  }

  return `${baseUrl}/assets/${video.id}/playlist.m3u8`;
}
