import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import * as TweetContainer from "containers/TweetContainer";

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
  tweet,
  user,
  retweetUser = null,
  retweets,
  reactions,

  userReactions = [],
  userRetweeted,
}) {
  const authContext = AuthContainer.useContext();
  const tweetContext = TweetContainer.useContext();

  const handleRetweet = async () => {
    if (userRetweeted) {
      tweetContext.removeRetweet(tweet.id);
    } else {
      tweetContext.retweet(tweet.id);
    }
  };

  const handleReact = r => {
    if (userReactions.includes(r)) {
      tweetContext.removeReaction(tweet.id, r);
    } else {
      tweetContext.addReaction(tweet.id, r);
    }
  };

  const handleDeleteTweet = () => {
    tweetContext.deleteTweet(tweet.id);
  };

  return (
    <div className="tweet">
      <img src="" alt={`${user.username}'s profile`} />
      <div className="tweet-content">
        {retweetUser && <div>Retweeted from: {retweetUser.username}</div>}
        <span>@{user.username} - 22h</span>
        <p>{tweet.text}</p>

        <div className="tweet-actions">
          <button
            onClick={handleRetweet}
            style={{ color: userRetweeted ? "green" : "black" }}
          >
            retweet {retweets}
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
                {REACTION_MAP[r]} {reactions[r] || 0}
              </button>
            ))}
          </div>
          <button>share</button>
          {user.id === authContext.loggedInUser.id && (
            <button onClick={handleDeleteTweet}>delete</button>
          )}
        </div>
      </div>
    </div>
  );
}
