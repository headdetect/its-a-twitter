import React from "react";

export const Context = React.createContext(null);

const DISPLAY_DURATION_MS = 6e3; // 10 seconds //

export function Provider({ children }) {
  const [alert, setAlert] = React.useState("");
  const [alertType, setAlertType] = React.useState("");
  const [alertIsVisible, setAlertIsVisible] = React.useState(false);
  const [alertTimeout, setAlertTimeout] = React.useState(null);

  const makeAlert = React.useCallback(
    (alertType, nodes) => {
      if (alertTimeout) {
        clearTimeout(alertTimeout);
      }

      setAlertType(alertType);
      setAlert(nodes);
      setAlertIsVisible(true);
      setAlertTimeout(
        setTimeout(() => setAlertIsVisible(false), DISPLAY_DURATION_MS),
      );
    },
    [alertTimeout],
  );

  return (
    <Context.Provider
      value={{
        // Actions //
        makeAlert,
      }}
    >
      <div
        className={`global-alert ${
          alertIsVisible ? "global-alert-visible" : ""
        }`}
      >
        <div className={`alert alert-${alertType}`}>{alert}</div>
      </div>

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
