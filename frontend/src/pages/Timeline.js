import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import * as TweetContainer from "containers/TweetContainer";

import Page from "components/Page";
import CraftTweet from "components/CraftTweet";
import Tweet from "components/Tweet";

import "./Timeline.css";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

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
      <div className="content what-is-this">
        <div className="alert alert-info">
          <h3>
            <FontAwesomeIcon icon="circle-question" fixedWidth /> What is this?
          </h3>

          <p>
            This is a fully-fledged micro application that emulates some of the
            same features as Twitter. <em>It&apos;s-a-Twitter</em> was developed
            for educational purposes only. <br />
            <br />
            To read more about it, visit our{" "}
            <a
              href="https://github.com/headdetect/its-a-twitter/blob/master/README.md"
              target="_blank"
              rel="noopener noreferrer"
            >
              documentation
            </a>{" "}
            on GitHub
          </p>
        </div>
      </div>

      {isLoggedIn && (
        <>
          <CraftTweet onTweet={tweetActions.submitTweet} />

          <div className="divider content">
            <span>Timeline</span>
          </div>
        </>
      )}

      <div className="timeline-stream  content">
        {timeline.length === 0 && <>It&apos;s pretty empty here...</>}

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
