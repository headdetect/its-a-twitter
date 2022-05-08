/**
 * @fileoverview Container that stores the current profile information. As well as
 * gives access to actions to interact with user profile endpoints with the API
 */
import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import { API_URL } from "consts";

export const Context = React.createContext(null);

export function Provider({ children, profileUsername = null }) {
  const { authenticatedFetch, loggedInUser } = AuthContainer.useContext();

  const [profileUser, setProfileUser] = React.useState(undefined);
  const [profileUserStatus, setProfileUserStatus] = React.useState("loading"); // loading | finished | error

  const followUser = React.useCallback(
    async username => {
      const response = await authenticatedFetch(
        `${API_URL}/user/profile/${username}/follow`,
        {
          method: "PUT",
        },
      );

      if (response.status === 401) {
        // TODO: Make the user log-in
        return;
      }

      if (!response.ok) {
        throw new Error("Error following user");
      }

      setProfileUser(oldProfileUser => ({
        ...oldProfileUser,
        followers: [...oldProfileUser.followers, loggedInUser],
      }));
    },
    [authenticatedFetch, loggedInUser],
  );
  const unfollowUser = React.useCallback(
    async username => {
      const response = await authenticatedFetch(
        `${API_URL}/user/profile/${username}/follow`,
        {
          method: "DELETE",
        },
      );

      if (response.status === 401) {
        // TODO: Make the user log-in
        return;
      }

      if (!response.ok) {
        throw new Error("Error unfollowing user");
      }

      setProfileUser(oldProfileUser => ({
        ...oldProfileUser,
        followers:
          oldProfileUser.followers?.filter(
            u => u.username !== loggedInUser.username,
          ) ?? [],
      }));
    },
    [authenticatedFetch, loggedInUser],
  );

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
        followUser,
        unfollowUser,
        setProfileUser,

        // State //
        profileUser,
        profileUserStatus,
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
