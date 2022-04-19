import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import * as UserContainer from "containers/UserContainer";
import * as TweetContainer from "containers/TweetContainer";

import TimelineStream from "components/tweet/TimelineStream";
import ProfileInfo from "components/user/ProfileInfo";
import Page from "components/Page";

export default function Profile({ username }) {
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
        <UserContainer.Provider profileUsername={username}>
          <Page title={`@${username}`}>
            <TweetContainer.Provider>
              <ProfileInfo />
              <TimelineStream />
            </TweetContainer.Provider>
          </Page>
        </UserContainer.Provider>
      </div>
    </AuthContainer.Provider>
  );
}
