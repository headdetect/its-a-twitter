/**
 * @fileoverview Container that stores the current user information if logged in
 * and gives the ability to log a user in otherwise.
 */
import React from "react";
import PropTypes from "prop-types";
import useHttpService from "../hooks/useHttpService";

export const Context = React.createContext(null);

const API_URL = process.env.REACT_APP_API_URL;

export function Provider(props) {
  const { authenticatedFetch, saveCredentials } = useHttpService();
  const [currentUser, setCurrentUser] = React.useState(null);

  React.useEffect(() => {
    const apiToken = localStorage.getItem("apiToken");

    if (apiToken) {
      saveCredentials(apiToken);
    }
  }, []);

  const logout = React.useCallback(() => {
    setCurrentUser(null);
    localStorage.removeItem("apiToken");
  }, []);

  const login = React.useCallback(async (username, password) => {
    if (currentUser) {
      // Log out first //
      logout();
    }

    let response;

    try {
      response = await fetch(`${API_URL}/login`, {
        method: "POST",
        body: JSON.stringify({
          username,
          password
        })
      });
    } catch (e) {
      console.error(e);
      throw new Error("There was a problem logging in. Try again.");
    }

    if (response.status === 401) {
      throw new Error("Invalid username or password.");
    }

    try {
      const { user, apiToken } = await response.json();

      if (!user || !apiToken) {
        // Bubbles up to this try's catch //
        throw new Error();
      }

      localStorage.setItem("apiToken", apiToken);
      setCurrentUser(user);
    } catch (e) {
      throw new Error("Server sent some weird stuff back. Try again.");
    }
  }, [currentUser]);

  const updateDisplayName = React.useCallback(async (displayName) => {
    try {
      const response = await authenticatedFetch(`${API_URL}/user/display-name`, {
        method: "PUT",
        body: JSON.stringify({
          displayName,
        })
      });

      const { user } = await response.json();

      setCurrentUser(user);
    } catch (e) {
      console.error(e);
      throw new Error("There was a problem updating the display name");
    }
  }, []);
 
  const updateProfilePic = React.useCallback(async (pic) => {
    try {
      const response = await authenticatedFetch(`${API_URL}/user/profile-pic`, {
        method: "PUT",
        body: pic
      });

      const { user } = await response.json();

      setCurrentUser(user);
    } catch (e) {
      console.error(e);
      throw new Error("There was a problem updating the display name");
    }
  }, []);

  return (
    <Context.Provider
      value={{
        currentUser,
        login,
        updateDisplayName,
        updateProfilePic,
      }}
    >
      {props.children}
    </Context.Provider>
  );
}

export function useContext() {
  const context = React.useContext(Context);

  if (!context) {
    throw new Error(
      "container can only be used in the context of a UserContainer.Provider"
    );
  }

  return context;
}

Provider.propTypes = {
  children: PropTypes.element.isRequired,
};