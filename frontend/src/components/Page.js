import React from "react";

import * as UserContainer from "containers/UserContainer";

import "./Page.css";
import LoginForm from "components/user/LoginForm";
import RegistrationForm from "components/user/RegistrationForm";

export default function Page({ children, title = "" }) {
  const { isLoggedIn, currentUser } = UserContainer.useContext();

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
          {title && <h1>{title}</h1>}

          {isLoggedIn ? (
            <>Hi {currentUser.username}</>
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
