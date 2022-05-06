import * as React from "react";

import * as AuthContainer from "containers/AuthContainer";
import FloatingLabelInput from "components/FloatingLabelInput";

import "./UserLoginForm.css"; // They'll have the same styles //

export default function UserRegistrationForm() {
  const { register } = AuthContainer.useContext();
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [email, setEmail] = React.useState("");

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

      if (!email) {
        setError("Email field is empty");
        return;
      }

      setError(null);
      setIsLoading(true);

      try {
        await register(email, username, password);

        // Redirect to own profile //
        window.location.href = `/profile/@${username}`;
      } catch (e) {
        setError(String(e));
      }

      setIsLoading(false);
    },
    [username, password, email, register],
  );

  return (
    <div className="profile-form">
      <div className="title">Register</div>

      <form onSubmit={handleSubmit}>
        <FloatingLabelInput
          label="Email"
          type="email"
          autoComplete="email"
          onChange={e => setEmail(e.target.value)}
        />

        <FloatingLabelInput
          label="Username"
          type="text"
          autoComplete="username"
          onChange={e => setUsername(e.target.value)}
        />

        <FloatingLabelInput
          label="Password"
          type="password"
          autoComplete="new-password"
          onChange={e => setPassword(e.target.value)}
        />

        {error && <div className="alert alert-danger">{error}</div>}

        <div className="profile-form-action">
          <button className="btn btn-light" type="submit" disabled={isLoading}>
            {isLoading ? "Registering..." : "Register"}
          </button>
        </div>
      </form>
    </div>
  );
}
