import React from "react";

export default function useHttpService() {
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

  return { authenticatedFetch, saveCredentials };
}
