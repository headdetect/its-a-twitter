import React from "react";
import { Link } from "react-router-dom";

import * as AuthContainer from "containers/AuthContainer";

import Page from "components/Page";

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
        <Page>
          This page does not exist.{" "}
          <Link to="/timeline">Click here to go back</Link>
        </Page>
      </div>
    </AuthContainer.Provider>
  );
}
