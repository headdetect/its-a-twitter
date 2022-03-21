import React from "react";

export default function useHttpService() {
  const authTokenRef = React.useRef();

  const authenticatedFetch = React.useCallback((url, options, ...rest) => {
    const { headers = {}, ...otherOptions } = options;

    const authenticatedOptions = {
      headers: {
        "Token": authTokenRef.current,
        ...headers,
      },
      ...otherOptions,
    };

    return fetch(url, authenticatedOptions, ...rest);
  }, []);

  const saveCredentials = React.useCallback(
    (authToken) => (authTokenRef.current = authToken), 
    []
  );

  return { authenticatedFetch, saveCredentials };
}