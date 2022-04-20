import React from "react";

export default function ProfileInfo({
  profileUser,
  isFollowingUser,
  canFollowUser,

  onFollowUser = _ => {},
  onUnfollowUser = _ => {},
}) {
  const handleChangeFollow = () => {
    if (isFollowingUser) {
      onUnfollowUser(profileUser.user.id);
    } else {
      onFollowUser(profileUser.user.id);
    }
  };

  return (
    <div>
      This is @{profileUser.user.username}.
      {canFollowUser && (
        <button onClick={handleChangeFollow}>
          {isFollowingUser ? "Unfollow" : "Follow"}
        </button>
      )}
    </div>
  );
}
