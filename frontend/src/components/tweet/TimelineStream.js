import React from "react";

import * as TweetContainer from "containers/TweetContainer";
import Tweet from "components/tweet/Tweet";

import "./TimelineStream.css";

export default function TimelineStream(props) {
  const { timelineStatus, timeline, timelineUsers } =
    TweetContainer.useContext();

  if (timelineStatus === "loading" || timelineStatus === "init") {
    return <>Loading...</>;
  }

  if (
    timelineStatus === "error" ||
    (timelineStatus === "finished" && !timeline)
  ) {
    return <>TODO: Error loading them hoes</>;
  }

  if (timelineStatus === "finished") {
    if (timeline.length === 0) {
      return <>There&apos;s nothing here :(</>;
    }

    return (
      <div className="timeline-stream">
        {timeline.map(timelineTweet => {
          const user = timelineUsers[timelineTweet.posterUserId];
          const retweetUser = timelineUsers[timelineTweet.retweeterUserId];
          const {
            tweet,
            reactionCount,
            retweetCount,
            userReactions,
            userRetweeted,
          } = timelineTweet;

          return (
            <Tweet
              key={tweet.id}
              tweet={tweet}
              user={user}
              retweetUser={retweetUser}
              retweets={retweetCount}
              reactions={reactionCount}
              userReactions={userReactions}
              userRetweeted={userRetweeted}
            />
          );
        })}
      </div>
    );
  }

  return <>Undefined state...</>;
}
