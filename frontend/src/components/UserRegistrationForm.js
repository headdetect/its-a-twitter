import * as React from "react";
import { useNavigate } from "react-router-dom";

import * as AuthContainer from "containers/AuthContainer";
import FloatingLabelInput from "components/FloatingLabelInput";

import "./UserLoginForm.css"; // They'll have the same styles //

const MAX_USERNAME_CHARS = 20;
const MIN_PASSWORD_CHARS = 8;

export default function UserRegistrationForm() {
  const navigate = useNavigate();
  const { register } = AuthContainer.useContext();
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [email, setEmail] = React.useState("");

  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState("");

  const [validationErrors, setValidationErrors] = React.useState({
    username: null,
    password: null,
  });

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

      let v = {
        username: null,
        password: null,
      };

      if (
        !username.match(/^[a-z0-9\_\-]*$/gi) ||
        username.length > MAX_USERNAME_CHARS
      ) {
        v.username = `Text must be ${MAX_USERNAME_CHARS} chars and only contain: letters, numbers, _, or -`;
      }

      if (password.length < MIN_PASSWORD_CHARS) {
        v.password = `Password must be at least ${MIN_PASSWORD_CHARS} chars`;
      }

      setError(null);

      if (v.username || v.password) {
        setValidationErrors(v);
        return;
      }

      setIsLoading(true);

      try {
        await register(email, username, password);

        // Redirect to own profile //
        navigate(`/profile/@${username}`);
        location.reload(); // Reload after //
      } catch (e) {
        setError(String(e));
      }

      setIsLoading(false);
    },
    [username, password, email, register, navigate],
  );

  return (
    <div className="profile-form">
      <div className="title">Register</div>

      <div className="alert alert-info what-is-this">
        This is a demo app. All data will be refreshed every 2 hours.
      </div>

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
          errorText={validationErrors.username}
          onChange={e => setUsername(e.target.value)}
        />

        <FloatingLabelInput
          label="Password"
          type="password"
          autoComplete="new-password"
          errorText={validationErrors.password}
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
