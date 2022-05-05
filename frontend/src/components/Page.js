import React, { forwardRef } from "react";

import * as AuthContainer from "containers/AuthContainer";
import UserLoginForm from "components/UserLoginForm";
import UserRegistrationForm from "components/UserRegistrationForm";

import "./Page.css";
import useClickedOutside from "hooks/useClickedOutside";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import UserAvatar from "components/UserAvatar";

function UserPanelPopover({ onClose, ...theRest }) {
  const ref = React.useRef();

  const { loggedInUser, logout } = AuthContainer.useContext();

  // Listen in this component so we unhook with unmount //
  useClickedOutside(ref, onClose);

  return (
    <div className="user-panel-popover" ref={ref} {...theRest}>
      <ul>
        <li className="section-header">Profile</li>
        <li>
          <a
            href={`/profile/@${loggedInUser.username}`}
            rel="noopener noreferrer"
          >
            <FontAwesomeIcon icon="user" />
            Your Profile
          </a>
        </li>
        <li>
          <a
            onClick={() => {
              logout(true);
            }}
            rel="noopener noreferrer"
            href=""
          >
            <FontAwesomeIcon icon="sign-out-alt" />
            Logout
          </a>
        </li>

        <li className="section-header">About</li>
        <li>
          <a href={"/about"} rel="noopener noreferrer">
            <FontAwesomeIcon icon="file-alt" />
            Docs
          </a>
        </li>
        <li>
          <a
            href="https://github.com/headdetect/its-a-twitter"
            target="_blank"
            rel="noopener noreferrer"
          >
            <FontAwesomeIcon icon={["fab", "github"]} />
            Github
          </a>
        </li>
      </ul>
    </div>
  );
}

export default function Page({ children, title = "" }) {
  const { isLoggedIn, loggedInUser } = AuthContainer.useContext();

  const [isUserPanelOpened, setIsUserPanelOpened] = React.useState(false);
  const [userPanelType, setUserPanelType] = React.useState(null);

  React.useEffect(() => {
    if (isLoggedIn) {
      setUserPanelType(null);
    }
  }, [isLoggedIn]);

  return (
    <>
      <div className="navbar">
        <div className="navbar-content">
          <div className="navbar-title">
            <a href="/timeline">
              <FontAwesomeIcon icon="kiwi-bird" />
            </a>

            {title}
          </div>

          {isLoggedIn ? (
            <div className="user-profile">
              <button onClick={() => setIsUserPanelOpened(p => !p)}>
                <UserAvatar
                  user={loggedInUser}
                  size={25}
                  style={{ marginRight: "calc(var(--spacing) * 1.5)" }}
                />
                <span>@{loggedInUser.username}</span>
                <FontAwesomeIcon icon="caret-down" />
              </button>

              {isUserPanelOpened && (
                <UserPanelPopover onClose={() => setIsUserPanelOpened(false)} />
              )}
            </div>
          ) : (
            <div className="user-login-register">
              <button className="btn" onClick={() => setUserPanelType("login")}>
                Log in
              </button>
              <button
                className="btn"
                onClick={() => setUserPanelType("register")}
              >
                Register
              </button>
            </div>
          )}
        </div>
      </div>

      {userPanelType && (
        <div className="panel">
          {userPanelType === "login" && <UserLoginForm />}
          {userPanelType === "register" && <UserRegistrationForm />}
        </div>
      )}

      {children}
    </>
  );
}
