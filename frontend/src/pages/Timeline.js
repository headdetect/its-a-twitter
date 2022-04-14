import React from "react";

import * as UserContainer from "containers/UserContainer";
import * as TweetContainer from "containers/TweetContainer";

import TimelineStream from "components/tweet/TimelineStream";
import Page from "components/Page";
import CraftTweet from "components/tweet/CraftTweet";

export default function Timeline() {
  return (
    <UserContainer.Provider>
      <div
        style={{
          display: "flex",
          height: "100%",
          width: "100%",
        }}
      >
        <Page title="Timeline">
          <TweetContainer.Provider>
            <CraftTweet />

            <TimelineStream />
          </TweetContainer.Provider>
        </Page>
      </div>
    </UserContainer.Provider>
  );
}
