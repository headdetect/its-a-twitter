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
  const { isLoggedIn } = AuthContainer.useContext();
  const {
    setTimeline,
    fetchAndStoreTimeline,
    timelineStatus,
    timeline,
    timelineUsers,
    ...tweetActions
  } = TweetContainer.useContext();

  const [updatedAt, setUpdatedAt] = React.useState(new Date());

  React.useEffect(() => {
    fetchAndStoreTimeline();
  }, [fetchAndStoreTimeline]);

  React.useEffect(() => {
    const updateTimestampsInterval = setInterval(() => {
      setUpdatedAt(new Date());
    }, 30e3); // Every 30 seconds //

    return () => {
      clearInterval(updateTimestampsInterval);
    };
  });

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
      {isLoggedIn && <CraftTweet onTweet={tweetActions.submitTweet} />}

      <div className="timeline-stream">
        {timeline.length === 0 && <>There&apos;s nothing here :(</>}

        {timeline.map(timelineTweet => {
          const user = timelineUsers[timelineTweet.posterUserId];
          const retweetUser = timelineUsers[timelineTweet.retweeterUserId];

          const isRecent =
            Date.now() - timelineTweet.tweet.createdAt < 1e3 * 60 * 60; // 1hr //

          return (
            <Tweet
              key={timelineTweet.tweet.id + (isRecent ? updatedAt : 0)} // Force the most recent ones to re-render
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
