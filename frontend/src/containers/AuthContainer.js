/**
 * @fileoverview Container that stores the current user information if logged in
 * and gives the ability to log a user in otherwise.
 */
import React from "react";

export const Context = React.createContext(null);

export function Provider(props) {
  const authCredentials = React.useRef({
    username: null,
    authToken: null,
  });

  const authenticatedFetch = React.useCallback((url, options = {}) => {
    const { headers = {}, ...otherOptions } = options;

    const authenticatedOptions = {
      headers: {
        ...authCredentials.current,
        ...headers,
      },
      ...otherOptions,
    };

    return fetch(url, authenticatedOptions);
  }, []);

  const saveCredentials = React.useCallback((authToken, username) => {
    authCredentials.current = {
      authToken,
      username,
    };
  }, []);

  return (
    <Context.Provider
      value={{
        authenticatedFetch,
        saveCredentials,
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
      "container can only be used in the context of a UserContainer.Provider",
    );
  }

  return context;
}
