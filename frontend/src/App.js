import * as React from "react";

import "./App.css";

import * as ApiContainer from "containers/ApiContainer";

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
  faCircleQuestion,
  faClose,
} from "@fortawesome/free-solid-svg-icons";
import { faGithub } from "@fortawesome/free-brands-svg-icons";
import { Loading } from "pages/Loading";

const BASE_ROUTE = process.env.REACT_APP_BASE_ROUTE ?? "";
const LOAD_WAIT_MS = 2e3; // After 3 seconds //

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
  faCircleQuestion,
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
  const { isServerLoaded, serverError, updateServerIsLoaded } =
    ApiContainer.useContext();

  // Don't start rendering the load page until some x time has passed //
  const [isRenderingLoader, setIsRenderingLoader] = React.useState(false);

  React.useEffect(() => {
    updateServerIsLoaded();
  }, [updateServerIsLoaded]);

  React.useEffect(() => {
    const handle = setTimeout(() => setIsRenderingLoader(true), LOAD_WAIT_MS);
    return () => clearTimeout(handle);
  }, []);

  return (
    <div className="app">
      {/* Prevent rendering if the server is still coming up */}
      {isServerLoaded ? (
        <Route />
      ) : (
        <Loading error={serverError} visible={isRenderingLoader} />
      )}
    </div>
  );
}

export default App;
