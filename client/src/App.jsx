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
      <header className="App-header">
        <h1>Video Streaming App</h1>
        <VideoList videos={videos} onSelect={setSelectedVideoId} />
        { selectedVideoId && <VideoPlayer src={getVideoSrc(selectedVideoId)} /> }
      </header>
    </div>
  );
}

function VideoList({ videos, onSelect }) {
  return (
    <div>
      <h2>Available Videos</h2>
      <ul>
        {videos.map((id) => (
          <li key={id}>
            <button onClick={() => onSelect(id)}>{id}</button>
          </li>
        ))}
      </ul>
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
      <h2>Playing: {src}</h2>
      <video ref={ref} controls style={{ width: "300px" }} />
    </div>
  );
}

function getVideoSrc(videoId) {
  return `${baseUrl}/assets/${videoId}/playlist.m3u8`;
}
