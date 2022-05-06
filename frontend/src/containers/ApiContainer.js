import React from "react";
import { API_URL } from "consts";

export const Context = React.createContext(null);

const MAX_TIMEOUT_MS = 60e3; // A full minute //

export function Provider({ children }) {
  const [isServerLoaded, setIsServerLoaded] = React.useState(false);
  const [serverError, setServerError] = React.useState("");

  const updateServerIsLoaded = React.useCallback(async () => {
    const controller = new AbortController();
    const id = setTimeout(() => controller.abort(), MAX_TIMEOUT_MS);

    try {
      const response = await fetch(`${API_URL}/`, {
        signal: controller.signal,
      });

      clearTimeout(id);

      if (response.status !== 200) {
        // to be caught by the try/catch handler //
        throw new Error("Received error: " + response.statusText);
      }

      setIsServerLoaded(true);
    } catch (e) {
      setServerError(String(e));
    }
  }, []);

  return (
    <Context.Provider
      value={{
        // Actions //
        updateServerIsLoaded,

        // State //
        isServerLoaded,
        serverError,
      }}
    >
      {children}
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
