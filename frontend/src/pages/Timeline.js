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
  const authContext = AuthContainer.useContext();
  const tweetContext = TweetContainer.useContext();

  if (tweetContext.timelineStatus === "loading") {
    return <>Loading...</>;
  }

  if (
    tweetContext.timelineStatus === "error" ||
    (tweetContext.timelineStatus === "finished" && !tweetContext.timeline)
  ) {
    return <>TODO: Error loading them hoes</>;
  }

  return (
    <>
      <CraftTweet onTweet={tweetContext.tweet} />

      <div className="timeline-stream">
        {tweetContext.timeline.length === 0 && (
          <>There&apos;s nothing here :(</>
        )}

        {tweetContext.timeline.map(timelineTweet => {
          const user = tweetContext.timelineUsers[timelineTweet.posterUserId];
          const retweetUser =
            tweetContext.timelineUsers[timelineTweet.retweeterUserId];

          return (
            <Tweet
              key={timelineTweet.tweet.id}
              user={user}
              timelineTweet={timelineTweet}
              retweetUser={retweetUser}
              isOwnTweet={user.id === authContext.loggedInUser.id}
              // Actions //
              onRemoveRetweet={tweetContext.removeRetweet}
              onRetweet={tweetContext.retweet}
              onRemoveReaction={tweetContext.removeReaction}
              onReaction={tweetContext.addReaction}
              onDeleteTweet={tweetContext.deleteTweet}
            />
          );
        })}
      </div>
    </>
  );
}
