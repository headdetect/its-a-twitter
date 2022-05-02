import React, { forwardRef } from "react";

import * as AuthContainer from "containers/AuthContainer";
import LoginForm from "components/LoginForm";
import RegistrationForm from "components/RegistrationForm";

import "./Page.css";
import useClickedOutside from "hooks/useClickedOutside";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

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
              logout();
              window.location.href = "/";
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
                <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAACWCAAAAAAZai4+AAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QA/4ePzL8AAAAHdElNRQfmBBwKOCIs7mQRAAACaElEQVR42u3V3XKiMAAF4L7/G5xErYBaq9v6s6OtLlD/ocBDbQCB7Ey7FbgwF+fcSIA5fgYTHhIj83BvAFlkkUUWWWQZErLIIsuEkEUWWSaELLLIMiFkkUWWCSGLLLJMCFlkkWVCyCKLLBNCFllkmRCyyLoD6wXXgz2K+Okwfht05OAtvqm/Ycn3rM9O0bj9pzEa5cej6AZV05JvWeEYReMSIz9PqEYzyNUpeJOY/axqXPI1K5g/SZSNUyy0S4Cbfv6BCP5valPyNWuXT/F1NMCmurRCPz94xKqclB5+ZwcvsOOmJTewouPx6JaNHeyrS8+Y5wdzTMqTLsRRfXzkHw1LfmalORWNIfDx3BP2r3M6srDNT29hVTdPMYiTqH+dtIYl9Vjl0pbv2a/289M+utXNYRfrZJHaWpTUY20A2wsCz84ekcAhP32A0O52IV0hTu1KarH812W2u8Q2xkkiq0ap3z5VM7FqW1KHVcZLZ7xfzX9fv6geoxO3LWnEUrtNkAyrf+tQv3gQkJe2JY1YYdo4KbbFJabatchCF08tS+qwYutxW8y4mv91saQttfaqLOCcpb5jNimpNVtjWJ9ZtYPXJLkI7NLRDkJ7Zvt0ga3RCdqU1GOpLcfxgovnoJt2zNDzosTv6W9Z9QgX2VdOWpTUZCUbke+E3ewXRkNAqHfwSFt3c1jp8ld/e7d5SV1Wcp5ZUtqLMB/Fa0dKZ60Vqpfy7srrhU1LbmTdNWSRRZYJIYssskwIWWSRZULIIossE0IWWWSZELLIIsuEkEUWWSaELLLIMiFkkUWWCSGLLLJMCFlkkWVC/gI+i6YMw2buTAAAACV0RVh0ZGF0ZTpjcmVhdGUAMjAyMi0wNC0yOFQxOTo1NjozNCswOTowMA2Q9vkAAAAldEVYdGRhdGU6bW9kaWZ5ADIwMjItMDQtMjhUMTk6NTY6MzQrMDk6MDB8zU5FAAAAAElFTkSuQmCC" />
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
          {userPanelType === "login" && <LoginForm />}
          {userPanelType === "register" && <RegistrationForm />}
        </div>
      )}

      {children}
    </>
  );
}
