import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import LoginForm from "components/LoginForm";
import RegistrationForm from "components/RegistrationForm";

import "./Page.css";

export default function Page({ children, title = "" }) {
  const { isLoggedIn, loggedInUser, logout } = AuthContainer.useContext();

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
          <div>{title && <h1>{title}</h1>}</div>

          {isLoggedIn ? (
            <div>
              Hi {loggedInUser.username}.{" "}
              <button onClick={() => logout()}>Log out?</button>
            </div>
          ) : (
            <div>
              <button onClick={() => setUserPanelType("login")}>Log in</button>
              <button onClick={() => setUserPanelType("register")}>
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
