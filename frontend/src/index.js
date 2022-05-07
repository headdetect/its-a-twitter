import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";

import * as ApiContainer from "containers/ApiContainer";
import * as AppContainer from "containers/AppContainer";

ReactDOM.render(
  <React.StrictMode>
    <ApiContainer.Provider>
      <AppContainer.Provider>
        <App />
      </AppContainer.Provider>
    </ApiContainer.Provider>
  </React.StrictMode>,
  document.getElementById("root"),
);
