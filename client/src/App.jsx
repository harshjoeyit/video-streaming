import { useRef, useEffect, useState } from "react";
import "./App.css";
import Hls from "hls.js";

const baseUrl = "http://localhost:9090";

export default function App() {
  // Call /assets API to list of available videos
  const [videos, setVideos] = useState([]);
  const [selectedVideoId, setSelectedVideoId] = useState(null);

  useEffect(() => {
    fetch(`${baseUrl}/assets`)
      .then((response) => response.json())
      .then((data) => {
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
        <VideoList videos={videos} onSelect={setSelectedVideoId} />
        {selectedVideoId && <VideoPlayer src={getVideoSrc(selectedVideoId)} />}
      </div>
    </div>
  );
}

function VideoList({ videos, onSelect }) {
  return (
    <div>
      <h4>Available Videos</h4>
      <ol>
        {videos.map((id) => (
          <li key={id} style={{ fontSize: "14px" }}>
            <button onClick={() => onSelect(id)}>{id}</button>
          </li>
        ))}
      </ol>
    </div>
  );
}

function VideoPlayer({ src }) {
  const ref = useRef(null);

  useEffect(() => {
    if (Hls.isSupported() && ref.current) {
      const hls = new Hls();
      hls.loadSource(src);
      console.log("HLS Source Loaded:", src);
      hls.attachMedia(ref.current);
      return () => hls.destroy();
    }
  }, [src]);

  return (
    <div>
      <h4>Playing: {src}</h4>
      <video ref={ref} controls style={{ width: "300px" }} />
    </div>
  );
}

function getVideoSrc(videoId) {
  return `${baseUrl}/assets/${videoId}/playlist.m3u8`;
}
