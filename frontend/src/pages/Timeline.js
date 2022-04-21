import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import * as TweetContainer from "containers/TweetContainer";

import Page from "components/Page";
import CraftTweet from "components/CraftTweet";
import Tweet from "components/Tweet";

import "./Timeline.css";

export function Presenter() {
  return (
    <AuthContainer.Provider>
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          height: "100%",
          width: "100%",
        }}
      >
        <Page title="Timeline">
          <TweetContainer.Provider>
            <Timeline />
          </TweetContainer.Provider>
        </Page>
      </div>
    </AuthContainer.Provider>
  );
}

function Timeline() {
  const {
    setTimeline,
    refreshTimeline,
    timelineStatus,
    timeline,
    timelineUsers,
    ...tweetActions
  } = TweetContainer.useContext();

  React.useEffect(() => {
    refreshTimeline();
  }, [refreshTimeline]);

  if (timelineStatus === "loading") {
    return <>Loading...</>;
  }

  if (
    timelineStatus === "error" ||
    (timelineStatus === "finished" && !timeline)
  ) {
    return <>TODO: Error loading them hoes</>;
  }

  return (
    <>
      <CraftTweet onTweet={tweetActions.tweet} />

      <div className="timeline-stream">
        {timeline.length === 0 && <>There&apos;s nothing here :(</>}

        {timeline.map(timelineTweet => {
          const user = timelineUsers[timelineTweet.posterUserId];
          const retweetUser = timelineUsers[timelineTweet.retweeterUserId];

          return (
            <Tweet
              key={timelineTweet.tweet.id}
              user={user}
              timelineTweet={timelineTweet}
              retweetUser={retweetUser}
              // Actions //
              onRemoveRetweet={tweetActions.removeRetweet}
              onRetweet={tweetActions.retweet}
              onRemoveReaction={tweetActions.removeReaction}
              onReaction={tweetActions.addReaction}
              onDeleteTweet={tweetActions.deleteTweet}
            />
          );
        })}
      </div>
    </>
  );
}
