import React from "react";

import * as UserContainer from "containers/UserContainer";
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
  const userContext = UserContainer.useContext();
  const tweetContext = TweetContainer.useContext();

  const [retweetStatus, setRetweetStatus] = React.useState(
    userRetweeted ? "retweeted" : "default", // default | retweeting | retweeted
  );

  const [reactionStatus, setReactionStatus] = React.useState(
    ALLOWED_REACTIONS.reduce((acc, curr) => {
      acc[curr] = userReactions.includes(curr) ? "reacted" : "default"; // default | reacting | reacted
      return acc;
    }, {}),
  );

  // We subtract one if this is a previous retweet given that the count
  // already includes our retweet in the numbers.
  // This allows us to track that ±1 in the state.
  const adjustedRetweetCount = userRetweeted ? retweets - 1 : retweets;
  const adjustedReactionCount = Object.keys(reactions).reduce((acc, curr) => {
    if (userReactions.includes(curr)) {
      acc[curr] -= 1;
    }

    return acc;
  }, reactions);

  const handleRetweet = () => {
    if (retweetStatus === "retweeting") {
      // prevent change in state while pending //
      return;
    }

    if (retweetStatus === "default") {
      tweetContext
        .retweet(tweet.id)
        .then(() => setRetweetStatus("retweeted"))
        .catch(() => setRetweetStatus("default")); // reset it back //

      setRetweetStatus("retweeting");
    } else if (retweetStatus === "retweeted") {
      tweetContext
        .removeRetweet(tweet.id)
        .then(() => setRetweetStatus("default"))
        .catch(() => setRetweetStatus("retweeted")); // reset it back //

      setRetweetStatus("retweeting");
    } else {
      throw new Error("Retweet status has undefined state");
    }
  };

  const handleReact = r => {
    if (reactionStatus[r] === "reacting") {
      // prevent change in state while pending //
      return;
    }

    const update = val => {
      setReactionStatus(a => ({
        ...a,
        [r]: val,
      }));
    };

    if (reactionStatus[r] === "default") {
      tweetContext
        .addReaction(tweet.id, r)
        .then(() => update("reacted"))
        .catch(() => update("default")); // reset it back //

      update("reacting");
    } else if (reactionStatus[r] === "reacted") {
      tweetContext
        .removeReaction(tweet.id, r)
        .then(() => update("default"))
        .catch(() => update("reacted")); // reset it back //

      update("reacting");
    } else {
      throw new Error("Reaction status has undefined state");
    }
  };

  return (
    <div
      className="tweet"
      style={{ borderColor: retweetUser ? "red" : "blue" }}
    >
      <img src="" alt={`${user.username}'s profile`} />
      <div className="tweet-content">
        <span>@{user.username} - 22h</span>
        <p>{tweet.text}</p>

        <div className="tweet-actions">
          <button
            onClick={handleRetweet}
            style={{ color: retweetStatus != "default" ? "green" : "black" }}
          >
            retweet {adjustedRetweetCount + (retweetStatus !== "default" ? 1 : 0)}
          </button>
          <div>
            {ALLOWED_REACTIONS.map(r => (
              <button
                key={r}
                onClick={() => handleReact(r)}
                style={{
                  color: reactionStatus[r] !== "default" ? "green" : "black",
                }}
              >
                {REACTION_MAP[r]}{" "}
                {(adjustedReactionCount.hasOwnProperty(r)
                  ? adjustedReactionCount[r]
                  : 0) + (reactionStatus[r] !== "default" ? 1 : 0)}
              </button>
            ))}
          </div>
          <button>share</button>
        </div>
      </div>
    </div>
  );
}
