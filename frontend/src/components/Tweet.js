import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import * as TimeUtils from "utils/TimeUtils";

import "./Tweet.css";

const API_URL = process.env.REACT_APP_API_URL;

const REACTION_MAP = {
  clap: "👏",
  party: "🎉",
  sad: "😔",
  heart: "❤️",
  thumbsup: "👍",
  thumbsdown: "👎",
};

const ALLOWED_REACTIONS = Object.keys(REACTION_MAP);

export default function Tweet({
  timelineTweet,
  user,
  retweetUser = null,

  onRemoveRetweet = _ => {},
  onRetweet = _ => {},
  onRemoveReaction = (_, __) => {},
  onReaction = (_, __) => {},
  onDeleteTweet = _ => {},
}) {
  const { loggedInUser, isLoggedIn } = AuthContainer.useContext();

  const {
    tweet,
    reactionCount = {},
    retweetCount = 0,
    userReactions = [],
    userRetweeted = false,
  } = timelineTweet;

  const handleRetweet = async () => {
    if (userRetweeted) {
      onRemoveRetweet(tweet.id);
    } else {
      onRetweet(tweet.id);
    }
  };

  const handleReact = r => {
    if (userReactions.includes(r)) {
      onRemoveReaction(tweet.id, r);
    } else {
      onReaction(tweet.id, r);
    }
  };

  const handleDeleteTweet = () => {
    onDeleteTweet(tweet.id);
  };

  const isOwnTweet = loggedInUser && loggedInUser.id === user.id;

  return (
    <div className="tweet">
      <img src="" alt={`${user.username}'s profile`} />
      <div className="tweet-content">
        {retweetUser && (
          <div>
            Retweeted from:{" "}
            <a
              href={`/profile/@${retweetUser.username}`}
              rel="noopener noreferrer"
            >
              @{retweetUser.username}
            </a>
          </div>
        )}
        <span>
          {" "}
          <a href={`/profile/@${user.username}`} rel="noopener noreferrer">
            @{user.username}
          </a>{" "}
          - {TimeUtils.toAgoString(new Date(tweet.createdAt * 1000))}
        </span>
        <p>{tweet.text}</p>

        {tweet.mediaPath && (
          <img src={`${API_URL}/asset/${tweet.mediaPath}`} alt="tweet media" />
        )}

        <div className="tweet-actions">
          <button
            onClick={handleRetweet}
            style={{ color: userRetweeted ? "green" : "black" }}
            disabled={!isLoggedIn || isOwnTweet}
          >
            retweet {retweetCount}
          </button>
          <div>
            {ALLOWED_REACTIONS.map(r => (
              <button
                key={r}
                onClick={() => handleReact(r)}
                style={{
                  color: userReactions.includes(r) ? "green" : "black",
                }}
                disabled={!isLoggedIn || isOwnTweet}
              >
                {REACTION_MAP[r]} {reactionCount[r] || 0}
              </button>
            ))}
          </div>
          <button>share</button>
          {isOwnTweet && <button onClick={handleDeleteTweet}>Delete</button>}
        </div>
      </div>
    </div>
  );
}
