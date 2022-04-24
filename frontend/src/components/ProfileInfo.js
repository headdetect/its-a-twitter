import React from "react";

import * as AuthContainer from "containers/AuthContainer";

export default function ProfileInfo({
  profileUser,

  onFollowUser = _ => {},
  onUnfollowUser = _ => {},
}) {
  const { isLoggedIn, loggedInUser } = AuthContainer.useContext();

  const isOwnProfile = profileUser.user.id === loggedInUser?.id;

  const profileFollowingLoggedInUser = React.useMemo(
    () =>
      profileUser?.following
        ? Boolean(profileUser.following.find(f => f.id === loggedInUser.id))
        : false,
    [profileUser, loggedInUser],
  );

  const loggedInUserFollowingProfile = React.useMemo(
    () =>
      profileUser?.followers
        ? Boolean(profileUser.followers.find(f => f.id === loggedInUser.id))
        : false,
    [profileUser, loggedInUser],
  );

  const handleChangeFollow = () => {
    if (loggedInUserFollowingProfile) {
      onUnfollowUser(profileUser.user.username);
    } else {
      onFollowUser(profileUser.user.username);
    }
  };

  return (
    <div>
      This is @{profileUser.user.username}.
      {isOwnProfile ? (
        <p>This is your profile</p>
      ) : (
        <button onClick={handleChangeFollow}>
          {loggedInUserFollowingProfile ? "Unfollow" : "Follow"}
        </button>
      )}
      {profileFollowingLoggedInUser && (
        <>@{profileUser.user.username} follows you</>
      )}
    </div>
  );
}
