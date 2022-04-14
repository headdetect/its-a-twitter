/**
 * @fileoverview Container that stores the current user information if logged in
 * and gives the ability to log a user in otherwise.
 */
import React from "react";
import PropTypes from "prop-types";
import useHttpService from "hooks/useHttpService";
import useEffectOnce from "hooks/useEffectOnce";

export const Context = React.createContext(null);

const API_URL = process.env.REACT_APP_API_URL;

export function Provider(props) {
  const { authenticatedFetch, saveCredentials, clearCredentials } = useHttpService();
  const [currentUser, setCurrentUser] = React.useState(null);
  const [isFinishedInitializing, setIsFinishedInitializing] =
    React.useState(false);

  const logout = React.useCallback(() => {
    setCurrentUser(null);
    localStorage.removeItem("authToken");
    localStorage.removeItem("username");
  }, []);

  const login = React.useCallback(
    async (username, password) => {
      if (currentUser) {
        // Log out first //
        logout();
      }

      let response;

      try {
        response = await fetch(`${API_URL}/user/login`, {
          method: "POST",
          body: JSON.stringify({
            username,
            password,
          }),
        });
      } catch (e) {
        console.error(e);
        throw new Error("There was a problem logging in. Try again.");
      }

      if (response.status === 401) {
        throw new Error("Invalid username or password.");
      }

      try {
        const { user, authToken } = await response.json();

        if (!user || !authToken) {
          // Bubbles up to this try's catch //
          throw new Error();
        }

        localStorage.setItem("authToken", authToken);
        localStorage.setItem("username", user.username);
        setCurrentUser(user);
      } catch (e) {
        throw new Error("Server sent some weird stuff back. Try again.");
      }
    },
    [currentUser, logout],
  );

  const updateProfilePic = React.useCallback(async pic => {
    // TODO
  }, []);

  const getOwnUser = React.useCallback(async () => {
    try {
      const response = await authenticatedFetch(`${API_URL}/user/self`);

      if (response.status === 401) {
        throw new Error("Not logged in.");
      }

      const obj = await response.json();

      return obj;
    } catch (e) {
      console.error(e);
      throw new Error("There was a problem logging in. Try again.");
    }
  }, [authenticatedFetch]);

  useEffectOnce(() => {
    const authToken = localStorage.getItem("authToken");
    const username = localStorage.getItem("username");

    if (authToken && username) {
      saveCredentials(authToken, username);

      getOwnUser()
        .then(user => {
          if (user?.user) {
            setCurrentUser(user.user);
          } else {
            logout();
          }
        })
        .catch(logout)
        .finally(() => setIsFinishedInitializing(true));
    } else {
      setIsFinishedInitializing(true);
    }
  }, [logout, getOwnUser, saveCredentials, setIsFinishedInitializing]);

  return (
    <Context.Provider
      value={{
        isLoggedIn: currentUser !== null,
        currentUser,
        login,
        updateProfilePic,
      }}
    >
      {isFinishedInitializing && props.children}
    </Context.Provider>
  );
}

export function useContext() {
  const context = React.useContext(Context);

  if (!context) {
    throw new Error(
      "container can only be used in the context of a UserContainer.Provider",
    );
  }

  return context;
}
