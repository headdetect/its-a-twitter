import * as React from "react";

import * as AuthContainer from "containers/AuthContainer";

export default function LoginForm() {
  const { login } = AuthContainer.useContext();
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState(null);

  const handleSubmit = React.useCallback(
    async e => {
      e.preventDefault();

      if (!username) {
        setError("Username field is empty");
        return;
      }

      if (!password) {
        setError("Password field is empty");
        return;
      }

      setError(null);
      setIsLoading(true);

      try {
        await login(username, password);
      } catch (e) {
        setError(String(e));
      }

      setIsLoading(false);
    },
    [username, password, login],
  );

  return (
    <div>
      <div>Log in.</div>
      {error && <div className="alert-error">{error}</div>}
      <form onSubmit={handleSubmit}>
        <input
          id="username"
          autoComplete="username"
          onChange={e => setUsername(e.target.value)}
        />
        <input
          id="password"
          type="password"
          autoComplete="current-password"
          onChange={e => setPassword(e.target.value)}
        />
        <button type="submit" disabled={isLoading}>
          {isLoading ? "Logging in..." : "Login"}
        </button>
      </form>
    </div>
  );
}
