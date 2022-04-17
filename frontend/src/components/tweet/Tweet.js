import React from "react";

import * as UserContainer from "containers/UserContainer";
import * as TweetContainer from "containers/TweetContainer";

import "./Tweet.css";

export default function Tweet({
  tweet,
  user,
  retweetUser = null,
  retweets,
  reactions,
}) {
  const userContext = UserContainer.useContext();
  const tweetContext = TweetContainer.useContext();

  return (
    <div
      className="tweet"
      style={{ borderColor: retweetUser ? "red" : "default" }}
    >
      <img src="" alt={`${user.username}'s profile`} />
      <div className="tweet-content">
        <span>@{user.username} - 22h</span>
        <p>{tweet.text}</p>

        <div className="tweet-actions">
          retweet {retweets}, react {reactions}, share
        </div>
      </div>
    </div>
  );
}
