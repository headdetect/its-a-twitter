import React from "react";

import "./Tweet.css";

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
  retweetUser,
  isOwnTweet,

  onRemoveRetweet = _ => {},
  onRetweet = _ => {},
  onRemoveReaction = (_, __) => {},
  onReaction = (_, __) => {},
  onDeleteTweet = _ => {},
}) {
  const { tweet, reactionCount, retweetCount, userReactions, userRetweeted } =
    timelineTweet;

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
          - 22h
        </span>
        <p>{tweet.text}</p>

        <div className="tweet-actions">
          <button
            onClick={handleRetweet}
            style={{ color: userRetweeted ? "green" : "black" }}
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
              >
                {REACTION_MAP[r]} {reactionCount[r] || 0}
              </button>
            ))}
          </div>
          <button>share</button>
          {isOwnTweet && <button onClick={handleDeleteTweet}>delete</button>}
        </div>
      </div>
    </div>
  );
}
