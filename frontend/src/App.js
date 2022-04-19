import * as React from "react";

import "./App.css";

import NotFound from "./pages/NotFound";
import Profile from "./pages/Profile";
import Timeline from "./pages/Timeline";

function Route() {
  const path = window.location.pathname;
  const segments = path.split("/").filter(p => Boolean(p));

  if (path === "/") {
    location.pathname = "/timeline";
  }

  if (path === "/timeline") {
    return <Timeline />;
  }

  if (path.startsWith("/profile") && segments.length >= 2) {
    const [_, atUsername] = segments;

    if (!atUsername.startsWith("@")) {
      return <NotFound />;
    }

    const username = atUsername.substring(1);

    return <Profile username={username} />;
  }

  return <NotFound />;
}

function App() {
  return (
    <div className="app">
      <Route />
    </div>
  );
}

export default App;
