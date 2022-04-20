import React from "react";

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
        <Page>Idk how you got here bruv</Page>
      </div>
    </AuthContainer.Provider>
  );
}
