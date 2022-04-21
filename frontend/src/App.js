import * as React from "react";

import "./App.css";

import * as NotFound from "./pages/NotFound";
import * as Profile from "./pages/Profile";
import * as Timeline from "./pages/Timeline";
import * as SingleTweet from "./pages/SingleTweet";

function Route() {
  const path = window.location.pathname;
  const segments = path.split("/").filter(p => Boolean(p));

  if (path === "/") {
    location.pathname = "/timeline";
  }

  if (path === "/timeline") {
    return <Timeline.Presenter />;
  }

  if (path.startsWith("/profile") && segments.length >= 2) {
    const [_, atUsername] = segments;

    if (!atUsername.startsWith("@")) {
      return <NotFound.Presenter />;
    }

    const username = atUsername.substring(1);

    return <Profile.Presenter username={username} />;
  }

  if (path.startsWith("/tweet") && segments.length >= 2) {
    const [_, tweetId] = segments;

    return <SingleTweet.Presenter tweetId={tweetId} />;
  }

  return <NotFound.Presenter />;
}

function App() {
  return (
    <div className="app">
      <Route />
    </div>
  );
}

export default App;
