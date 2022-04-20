/**
 * @fileoverview Container that stores the current user information if logged in
 * and gives the ability to log a user in otherwise.
 */
import React from "react";

import * as AuthContainer from "containers/AuthContainer";

export const Context = React.createContext(null);

const API_URL = process.env.REACT_APP_API_URL;

export function Provider({ children, profileUsername = null }) {
  const { authenticatedFetch, saveCredentials, loggedInUser } =
    AuthContainer.useContext();

  const [profileUser, setProfileUser] = React.useState(undefined);
  const [profileUserStatus, setProfileUserStatus] = React.useState("loading"); // loading | finished | error

  const [loggedInUserFollowing, setLoggedInUserFollowing] = React.useState([]);

  const followingUser = React.useMemo(
    () =>
      profileUser
        ? Boolean(loggedInUserFollowing.find(f => f.id === profileUser.id))
        : false,
    [loggedInUserFollowing, profileUser],
  );

  const updateProfilePic = React.useCallback(async pic => {}, []);
  const followUser = React.useCallback(async username => {}, []);
  const unfollowUser = React.useCallback(async username => {}, []);

  const loadProfileUser = React.useCallback(
    async username => {
      try {
        const response = await authenticatedFetch(
          `${API_URL}/user/profile/${username}`,
        );

        if (response.status === 401) {
          setProfileUserStatus("error");
          return;
        }

        const obj = await response.json();

        if (obj.user) {
          setProfileUser(obj);
          setProfileUserStatus("finished");
        } else {
          setProfileUserStatus("error");
        }
      } catch (e) {
        setProfileUserStatus("error");
      }
    },
    [authenticatedFetch],
  );

  React.useEffect(() => {
    if (profileUsername) {
      loadProfileUser(profileUsername);
    }

    // TODO: Get all of the logged in user's followers
  }, [profileUsername, loadProfileUser]);

  return (
    <Context.Provider
      value={{
        // Actions //
        updateProfilePic,
        followUser,
        unfollowUser,

        // State //
        profileUser,
        profileUserStatus,
        followingUser,
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
