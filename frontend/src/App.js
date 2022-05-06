import * as React from "react";

import "./App.css";

import * as NotFound from "./pages/NotFound";
import * as Profile from "./pages/Profile";
import * as Timeline from "./pages/Timeline";
import * as SingleTweet from "./pages/SingleTweet";

import { library } from "@fortawesome/fontawesome-svg-core";
import {
  faKiwiBird,
  faCaretDown,
  faUser,
  faImage,
  faTrash,
  faSignInAlt,
  faSignOutAlt,
  faFileAlt,
  faRetweet,
  faCalendar,
  faClose,
} from "@fortawesome/free-solid-svg-icons";
import { faGithub } from "@fortawesome/free-brands-svg-icons";

const BASE_ROUTE = process.env.REACT_APP_BASE_ROUTE ?? "";

library.add(
  faKiwiBird,
  faCaretDown,
  faUser,
  faImage,
  faTrash,
  faSignInAlt,
  faSignOutAlt,
  faFileAlt,
  faRetweet,
  faCalendar,
  faClose,

  faGithub,
);

function Route() {
  const path = window.location.pathname;
  const segments = path
    .split("/")
    .filter(p => Boolean(p) && p != BASE_ROUTE.substring(1));

  const [root = "", secondary = null] = segments;

  const routes = {
    "": () => {
      location.pathname = "/timeline";

      return null;
    },
    timeline: () => <Timeline.Presenter />,
    profile: () => {
      if (!secondary || !secondary.startsWith("@")) {
        return <NotFound.Presenter />;
      }

      const username = secondary.substring(1);

      return <Profile.Presenter username={username} />;
    },
    tweet: () => {
      if (!secondary) {
        return <NotFound.Presenter />;
      }

      return <SingleTweet.Presenter tweetId={+secondary} />;
    },
  };

  const component = routes[root];

  if (component) {
    const result = component();

    if (result) {
      return result;
    }
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
