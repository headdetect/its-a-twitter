import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import * as UserContainer from "containers/UserContainer";

export default function RegistrationForm(props) {
  const {
    profileUserStatus,
    profileUser,
    followUser,
    unfollowUser,
    followingUser,
  } = UserContainer.useContext();

  if (profileUserStatus === "loading") {
    return <>Loading...</>;
  }

  if (profileUserStatus === "error") {
    return <>Error fetching user</>;
  }

  if (profileUserStatus !== "finished") {
    return <>Unknown profile state specified: {profileUserStatus}</>;
  }

  const handleChangeFollow = () => {
    if (followingUser) {
      unfollowUser(profileUser.id);
    } else {
      followUser(profileUser.id);
    }
  };

  return (
    <div>
      This is @{profileUser.username}.
      <button onClick={handleChangeFollow}>
        {followingUser ? "Unfollow" : "Follow"}
      </button>
    </div>
  );
}
