import React from "react";

import "./Tweet.css";
import "./CraftTweet.css";

const ACCEPTABLE_MIME_TYPES = ["image/jpeg", "image/gif", "image/png"];

export default function CraftTweet({ onTweet = async (_, __) => {} }) {
  const [text, setText] = React.useState("");
  const [file, setFile] = React.useState(null);
  const [error, setError] = React.useState(null);

  const fileInputRef = React.useRef(null);

  const handleSubmitTweet = async () => {
    if (file && !ACCEPTABLE_MIME_TYPES.includes(file.type)) {
      setError(
        "This file type is not allowed. You must choose an image type (gif, png, jpeg)",
      );
    }

    await onTweet(text.trim(), file);

    setText("");
    setFile(null);
    setError(null);

    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  return (
    <div className="craft-tweet-container">
      <div className="tweet craft-tweet" style={{ borderColor: "green" }}>
        <img src="" alt="Your avatar" />
        <div className="tweet-content">
          <div>
            <textarea
              onChange={e => setText(e.target.value)}
              placeholder="Whats up dude?"
              value={text}
            />
          </div>

          <div>
            <input
              type="file"
              onChange={e => setFile(e.target.files?.[0])}
              multiple={false}
              accept={ACCEPTABLE_MIME_TYPES.join(", ")}
              ref={fileInputRef}
            />
          </div>

          <button
            onClick={handleSubmitTweet}
            disabled={!file && text.trim() === ""}
          >
            Post
          </button>
          {error && (
            <div
              style={{
                background: "#f49697",
                border: "1px solid #e92e32",
                padding: 8,
              }}
            >
              {error}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
