import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import * as ProfileContainer from "containers/ProfileContainer";
import * as TweetContainer from "containers/TweetContainer";

import ProfileInfo from "components/ProfileInfo";
import Page from "components/Page";

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
  const profileContext = ProfileContainer.useContext();

  if (profileContext.profileUserStatus === "loading") {
    return <>Loading...</>;
  }

  if (profileContext.profileUserStatus === "error") {
    return <>Error fetching user</>;
  }

  if (profileContext.profileUserStatus !== "finished") {
    return (
      <>Unknown profile state specified: {profileContext.profileUserStatus}</>
    );
  }

  return (
    <>
      <ProfileInfo
        profileUser={profileContext.profileUserStatus}
        isFollowingUser={false}
        onFollowUser={profileContext.profileUserStatus}
        onUnfollowUser={profileContext.profileUserStatus}
      />
      TODO: Add dem tweets boy
    </>
  );
}
