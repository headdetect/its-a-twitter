import React from "react";

export default function ProfileInfo({
  profileUser,
  isFollowingUser,

  onFollowUser = _ => {},
  onUnfollowUser = _ => {},
}) {
  const handleChangeFollow = () => {
    if (isFollowingUser) {
      onUnfollowUser(profileUser.id);
    } else {
      onFollowUser(profileUser.id);
    }
  };

  return (
    <div>
      This is @{profileUser.username}.
      <button onClick={handleChangeFollow}>
        {isFollowingUser ? "Unfollow" : "Follow"}
      </button>
    </div>
  );
}
