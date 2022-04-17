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
        {timeline.map(tweet => (
          <Tweet key={tweet.id} tweet={tweet} retweets={0} reactions={0} />
        ))}
      </div>
    );
  }

  return <>Undefined state...</>;
}
