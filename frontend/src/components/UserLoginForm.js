import * as React from "react";

import * as AuthContainer from "containers/AuthContainer";
import FloatingLabelInput from "components/FloatingLabelInput";

import "./UserLoginForm.css";

export default function UserLoginForm() {
  const { login } = AuthContainer.useContext();
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState("");

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

        // Refresh the page //
        window.location.reload();
      } catch (e) {
        setError(String(e));
      }

      setIsLoading(false);
    },
    [username, password, login],
  );

  return (
    <div className="profile-form">
      <div className="title">Log In</div>

      <form onSubmit={handleSubmit}>
        <FloatingLabelInput
          label="Username"
          type="text"
          autoComplete="username"
          onChange={e => setUsername(e.target.value)}
        />

        <FloatingLabelInput
          label="Password"
          type="password"
          autoComplete="current-password"
          onChange={e => setPassword(e.target.value)}
        />

        {error && <div className="alert alert-danger">{error}</div>}

        <div className="profile-form-action">
          <button className="btn btn-light" type="submit" disabled={isLoading}>
            {isLoading ? "Logging in..." : "Login"}
          </button>
        </div>
      </form>
    </div>
  );
}
