import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import * as ProfileContainer from "containers/ProfileContainer";
import * as TweetContainer from "containers/TweetContainer";

import ProfileInfo from "components/ProfileInfo";
import Page from "components/Page";

import "./Timeline.css";
import "./Profile.css";
import Tweet from "components/Tweet";

export function Presenter({ username }) {
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
        <Page title={`@${username}`}>
          <ProfileContainer.Provider profileUsername={username}>
            <TweetContainer.Provider>
              <Profile />
            </TweetContainer.Provider>
          </ProfileContainer.Provider>
        </Page>
      </div>
    </AuthContainer.Provider>
  );
}

function Profile() {
  const { updateProfilePic, loggedInUser } = AuthContainer.useContext();

  const { timeline, setTimeline, ...tweetActions } =
    TweetContainer.useContext(); // Just using the actions //

  const {
    profileUser,
    setProfileUser,
    profileUserStatus,
    followUser,
    unfollowUser,
  } = ProfileContainer.useContext();

  //  If we're on our own profile, update if logged in info changes  //
  React.useEffect(() => {
    if (loggedInUser?.username === profileUser?.user?.username) {
      setProfileUser(p => ({
        ...p,
        user: loggedInUser,
      }));
    }
  }, [loggedInUser, profileUser?.user?.username, setProfileUser]);

  React.useEffect(() => {
    if (!profileUser) {
      return;
    }

    setTimeline(profileUser.timeline?.tweets);
  }, [setTimeline, profileUser]);

  if (profileUserStatus === "loading") {
    return <>Loading...</>;
  }

  if (profileUserStatus === "error") {
    return <>Error fetching user</>;
  }

  if (profileUserStatus !== "finished") {
    return <>Unknown profile state specified: {profileUserStatus}</>;
  }

  if (!profileUser) {
    return <>Cannot find profile</>;
  }

  return (
    <>
      <ProfileInfo
        profileUser={profileUser}
        onFollowUser={followUser}
        onUnfollowUser={unfollowUser}
        onProfileImageChanged={updateProfilePic}
      />

      <div className="timeline-stream">
        {timeline.length === 0 && <>It&apos;s pretty empty here...</>}

        {timeline.map(timelineTweet => {
          // N.b If the originalTweetUser is not null, then this user is a retweeter
          const originalTweetUser =
            profileUser.timeline.users[timelineTweet.posterUserId];

          return (
            <Tweet
              key={timelineTweet.tweet.id}
              user={originalTweetUser ?? profileUser.user}
              timelineTweet={timelineTweet}
              retweetUser={originalTweetUser ? profileUser.user : null}
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
