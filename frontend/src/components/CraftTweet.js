import React from "react";

import "./Tweet.css";
import "./CraftTweet.css";

export default function CraftTweet({ onTweet = async (_, __) => {} }) {
  const [text, setText] = React.useState("");
  const [file, setFile] = React.useState(null);

  const fileInputRef = React.useRef(null);

  const handleSubmitTweet = async () => {
    await onTweet(text.trim(), file);

    setText("");
    setFile(null);

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
              ref={fileInputRef}
            />
          </div>

          <button
            onClick={handleSubmitTweet}
            disabled={!file && text.trim() === ""}
          >
            Post
          </button>
        </div>
      </div>
    </div>
  );
}
