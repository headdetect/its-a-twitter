import React from "react";

import "./Loading.css";

export function Loading({ error, visible }) {
  return (
    <>
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          height: "100%",
          width: "100%",
        }}
      >
        {visible && (
          <div className="loading-page">
            <h2>Loading...</h2>
            <p className="content">
              This may take a second while the API wakes up. <br />
              <a
                href="https://devcenter.heroku.com/articles/free-dyno-hours#dyno-sleeping"
                target="_blank"
                rel="noopener noreferrer"
              >
                Why?
              </a>
            </p>

            {error && (
              <div className="alert alert-danger">
                {error}. <br />
                Refresh the page in a few minutes
              </div>
            )}
          </div>
        )}
      </div>
    </>
  );
}
