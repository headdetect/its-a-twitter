import React from "react";

import "./Tweet.css";
import "./CraftTweet.css";

export default function CraftTweet({ onTweet = async _ => {} }) {
  const [text, setText] = React.useState("");

  const handleSubmitTweet = async () => {
    await onTweet(text);
    setText("");
  };

  return (
    <div className="craft-tweet-container">
      <div className="tweet craft-tweet" style={{ borderColor: "green" }}>
        <img src="" alt="Your avatar" />
        <div className="tweet-content">
          <textarea
            onChange={e => setText(e.target.value)}
            placeholder="Whats up dude?"
            value={text}
          />

          <button onClick={handleSubmitTweet}>Post</button>
        </div>
      </div>
    </div>
  );
}
