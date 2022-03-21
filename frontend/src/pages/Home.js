import * as React from "react";

export function Presenter() {
  return (
    <UserContainer.Provider>
      <TimelineContainer.Provider>
        <Home />
      </TimelineContainer.Provider>
    </UserContainer.Provider>
  )
}

function Login() {
  const userContext = UserContainer.useContext();

  React.useEffect(() => {
    if (userContext.currentUser) {
      // We're already logged in. Redirect to home //
    }

  }, [userContext.currentUser]);

  return (
    <>
      <div>Log in.</div>
      <form>
        <input id="username" />
        <input id="password" type="password" />
      </form>
    </>
  );
}