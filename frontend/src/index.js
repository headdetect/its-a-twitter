import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";

import * as ApiContainer from "containers/ApiContainer";

ReactDOM.render(
  <React.StrictMode>
    <ApiContainer.Provider>
      <App />
    </ApiContainer.Provider>
  </React.StrictMode>,
  document.getElementById("root"),
);
