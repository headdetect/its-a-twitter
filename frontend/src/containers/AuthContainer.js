/**
 * @fileoverview Container that stores the current user information if logged in
 * and gives the ability to log a user in otherwise.
 */
import useEffectOnce from "hooks/useEffectOnce";
import React from "react";

export const Context = React.createContext(null);

const API_URL = process.env.REACT_APP_API_URL;

export function Provider({ children }) {
  const [loggedInUser, setLoggedInUser] = React.useState(null);
  const [autoLoginStatus, setAutoLoginStatus] = React.useState("loading"); // loading | finished | error

  const authCredentials = React.useRef({
    username: null,
    authToken: null,
  });

  const authenticatedFetch = React.useCallback((url, options = {}) => {
    try {
      const { headers = {}, ...otherOptions } = options;

      const authenticatedOptions = {
        headers: {
          ...authCredentials.current,
          ...headers,
        },
        ...otherOptions,
      };

      return fetch(url, authenticatedOptions);
    } catch (e) {
      // Consume and rethrow for the caller to catch //
      throw e;
    }
  }, []);

  const saveCredentials = React.useCallback((authToken, username) => {
    authCredentials.current = {
      authToken,
      username,
    };
  }, []);

  const register = React.useCallback(
    async (email, username, password) => {
      let response;

      try {
        response = await fetch(`${API_URL}/user/register`, {
          method: "POST",
          body: JSON.stringify({
            email,
            username,
            password,
          }),
        });
      } catch (e) {
        console.error(e);
        throw new Error("There was a problem registering you. Try again.");
      }

      if (response.status === 400) {
        throw new Error("Error registering. Try again.");
      }

      if (response.status === 409) {
        throw new Error("Username already exists. Choose a different one.");
      }

      try {
        const { user, authToken } = await response.json();

        if (!user || !authToken) {
          // Bubbles up to this try's catch //
          throw new Error();
        }

        localStorage.setItem("authToken", authToken);
        localStorage.setItem("username", user.username);
        saveCredentials(authToken, user.username);
        setLoggedInUser(user);
      } catch (e) {
        throw new Error("Server sent some weird stuff back. Try again.");
      }
    },
    [saveCredentials, setLoggedInUser],
  );

  const logout = React.useCallback((forceRefresh = false) => {
    localStorage.removeItem("authToken");
    localStorage.removeItem("username");

    if (forceRefresh) {
      window.location.reload();
    } else {
      setLoggedInUser(null);
    }
  }, []);

  const login = React.useCallback(
    async (username, password) => {
      if (loggedInUser) {
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
        saveCredentials(authToken, user.username);
        setLoggedInUser(user);
      } catch (e) {
        throw new Error("Server sent some weird stuff back. Try again.");
      }
    },
    [loggedInUser, logout, saveCredentials],
  );

  const getOwnUser = React.useCallback(async () => {
    try {
      const response = await authenticatedFetch(`${API_URL}/user/self`);

      if (response.status === 401) {
        return null;
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
      setAutoLoginStatus("loading");

      saveCredentials(authToken, username);

      getOwnUser()
        .then(user => {
          if (user?.user) {
            setLoggedInUser(user.user);
          } else {
            logout();
          }
        })
        .catch(err => {
          setAutoLoginStatus("error");
        })
        .finally(() => setAutoLoginStatus("finished"));
    } else {
      setAutoLoginStatus("finished");
    }
  }, [logout, getOwnUser, saveCredentials, setAutoLoginStatus]);

  return (
    <Context.Provider
      value={{
        // Actions //
        authenticatedFetch,
        saveCredentials,
        login,
        register,
        logout,
        setLoggedInUser,

        // State //
        isLoggedIn: loggedInUser !== null,
        loggedInUser,
      }}
    >
      {/* Prevent rendering if we're not ready for attempted authenticated requests */}
      {autoLoginStatus !== "loading" && children}
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
