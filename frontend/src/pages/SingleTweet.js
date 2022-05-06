import React from "react";

import { useParams } from "react-router-dom";

import * as AuthContainer from "containers/AuthContainer";
import * as TweetContainer from "containers/TweetContainer";

import Page from "components/Page";
import Tweet from "components/Tweet";

import "./Timeline.css";

export function Presenter() {
  const { tweetId } = useParams();

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
        <Page title="Tweet">
          <TweetContainer.Provider>
            <SingleTweet tweetId={tweetId} />
          </TweetContainer.Provider>
        </Page>
      </div>
    </AuthContainer.Provider>
  );
}

function SingleTweet({ tweetId }) {
  const { isLoggedIn } = AuthContainer.useContext();
  const { getTweet, ...tweetActions } = TweetContainer.useContext();
  const [tweet, setTweet] = React.useState(null);
  const [tweetFetchStatus, setTweetFetchStatus] = React.useState("loading");

  React.useEffect(() => {
    getTweet(tweetId)
      .then(t => {
        setTweet(t.tweet);
        setTweetFetchStatus("finished");
      })
      .catch(() => {
        setTweetFetchStatus("error");
      });
  }, [getTweet, tweetId]);

  if (tweetFetchStatus === "loading") {
    return <>Loading...</>;
  }

  if (
    tweetFetchStatus === "error" ||
    (tweetFetchStatus === "finished" && !tweet)
  ) {
    return <>TODO: Error loading the tweet</>;
  }

  if (!tweet) {
    return <>Not sure how we got here</>;
  }

  return (
    <>
      <div className="timeline-stream">
        <Tweet
          key={tweet.id}
          user={tweet.tweet.user}
          timelineTweet={tweet}
          // Actions //
          onRemoveRetweet={isLoggedIn ? tweetActions.removeRetweet : undefined}
          onRetweet={isLoggedIn ? tweetActions.retweet : undefined}
          onRemoveReaction={
            isLoggedIn ? tweetActions.removeReaction : undefined
          }
          onReaction={isLoggedIn ? tweetActions.addReaction : undefined}
          onDeleteTweet={isLoggedIn ? tweetActions.deleteTweet : undefined}
        />
      </div>
    </>
  );
}
