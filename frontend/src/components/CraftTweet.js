import React from "react";

import "./Tweet.css";
import "./CraftTweet.css";
import * as AuthContainer from "containers/AuthContainer";
import UserAvatar from "components/UserAvatar";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const ACCEPTABLE_MIME_TYPES = ["image/jpeg", "image/gif", "image/png"];
const MAX_CHARS = 250;
const MAX_LINES = 4;

export default function CraftTweet({ onTweet = async (_, __) => {} }) {
  const { loggedInUser } = AuthContainer.useContext();

  const textAreaRef = React.useRef(null);

  const [text, setText] = React.useState("");
  const [mediaPreview, setMediaPreview] = React.useState("");
  const [file, setFile] = React.useState(null);
  const [error, setError] = React.useState(null);

  const fileInputRef = React.useRef(null);

  const updateTextAreaHeight = React.useCallback(() => {
    if (!textAreaRef.current) {
      return;
    }

    const borderPixels = 2;

    // Set height to 1px to get accurate scrollHeight value
    textAreaRef.current.style.height = "1px";

    const desiredHeight = textAreaRef.current.scrollHeight;
    const lineHeight = parseFloat(
      window
        .getComputedStyle(textAreaRef.current, null)
        .getPropertyValue("line-height"),
    );
    const maxHeight = lineHeight * MAX_LINES + borderPixels;
    const actualHeight = Math.min(maxHeight, desiredHeight);
    const adjustedHeight = actualHeight + borderPixels;

    textAreaRef.current.style.height = `${adjustedHeight}px`;
  }, []);

  // Run on mount //
  React.useEffect(() => {
    updateTextAreaHeight();
  }, [updateTextAreaHeight]);

  React.useLayoutEffect(() => {
    updateTextAreaHeight();
  }, [text, updateTextAreaHeight]);

  React.useEffect(() => {
    if (!fileInputRef.current) {
      return undefined;
    }

    // So we can access if this changes during unmount //
    const ref = fileInputRef.current;

    const handle = e => {
      const [file] = e.target.files;

      if (!ACCEPTABLE_MIME_TYPES.includes(file.type)) {
        setError("Only images are acceptable");
        return;
      }

      setError(null);

      // Convert to base64 so we can put in the image previewer //
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onloadend = () => {
        setMediaPreview(String(reader.result));
      };
    };

    ref.addEventListener("change", handle);

    return () => {
      ref?.removeEventListener("change", handle);
    };
  }); // no deps on purpose //

  const handleSubmitTweet = async () => {
    if (file && !ACCEPTABLE_MIME_TYPES.includes(file.type)) {
      setError(
        "This file type is not allowed. You must choose an image type (gif, png, jpeg)",
      );
    }

    await onTweet(text.trim(), file);

    setText("");
    setMediaPreview("");
    setFile(null);
    setError(null);

    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleUploadImage = () => {
    if (mediaPreview) {
      fileInputRef.current.value = "";
      setMediaPreview("");

      return;
    }

    fileInputRef.current.click();
  };

  let countModifier = "";
  if (text.trim().length === 0) {
    countModifier = "empty";
  } else if (text.trim().length > MAX_CHARS) {
    countModifier = "too-much";
  }

  return (
    <div className="craft-tweet-container">
      <div className="tweet craft-tweet" style={{ borderColor: "green" }}>
        <UserAvatar
          user={loggedInUser}
          size={50}
          style={{ marginRight: "calc(var(--spacing) * 2.5)" }}
        />

        <div className="tweet-content">
          <div className="tweet-editor">
            <textarea
              ref={textAreaRef}
              onChange={e => setText(e.target.value)}
              placeholder={`What's up @${loggedInUser.username}?`}
              value={text}
            />
            <span className={`count ${countModifier}`}>
              {text.length} / {MAX_CHARS}
            </span>
          </div>

          <div className="media-preview">
            {mediaPreview && <img src={mediaPreview} alt="uploaded media" />}
          </div>

          <div className="actions">
            <input
              type="file"
              style={{ display: "hidden" }}
              onChange={e => setFile(e.target.files?.[0])}
              multiple={false}
              accept={ACCEPTABLE_MIME_TYPES.join(", ")}
              ref={fileInputRef}
            />

            <button className="btn btn-upload" onClick={handleUploadImage}>
              <FontAwesomeIcon icon="image" />
              {mediaPreview ? "Remove Image" : "Add Image"}
            </button>

            <button
              className="btn btn-post"
              onClick={handleSubmitTweet}
              disabled={
                !file && (text.trim() === "" || text.trim().length > MAX_CHARS)
              }
            >
              Post
            </button>
          </div>

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
