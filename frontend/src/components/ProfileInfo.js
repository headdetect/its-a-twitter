import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import UserAvatar from "components/UserAvatar";

import "./ProfileInfo.css";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

const ACCEPTABLE_MIME_TYPES = ["image/jpeg", "image/gif", "image/png"];

export default function ProfileInfo({
  profileUser,

  onFollowUser = _ => {},
  onUnfollowUser = _ => {},
}) {
  const { isLoggedIn, loggedInUser } = AuthContainer.useContext();

  const fileInputRef = React.useRef(null);

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

  const handleProfileImageChanged = e => {
    const [file] = e.target.files;

    if (!ACCEPTABLE_MIME_TYPES.includes(file.type)) {
      // Silent rejection //
      return;
    }

    // TODO: upload and fetch new user
  };

  const handleChangeFollow = () => {
    if (loggedInUserFollowingProfile) {
      onUnfollowUser(profileUser.user.username);
    } else {
      onFollowUser(profileUser.user.username);
    }
  };

  const handleUploadProfileImage = () => {
    fileInputRef.current.click();
  };

  const joinedDate = new Date(profileUser.user.createdAt * 1000);

  const followingCount =
    profileUser.following?.length ?? profileUser.followingCount;
  const followersCount =
    profileUser.followers?.length ?? profileUser.followerCount;

  return (
    <div className="profile-info">
      <div className="profile-avatar">
        <input
          type="file"
          style={{ display: "hidden" }}
          onChange={handleProfileImageChanged}
          multiple={false}
          accept={ACCEPTABLE_MIME_TYPES.join(", ")}
          ref={fileInputRef}
        />

        <div className="upload-image">Replace</div>

        <UserAvatar
          user={profileUser.user}
          size={100}
          onClick={handleUploadProfileImage}
        />
      </div>

      <div className="profile-content">
        <div className="user">
          <span>@{profileUser.user.username}</span>
          {(isOwnProfile || profileFollowingLoggedInUser) && (
            <span className="subtitle">
              {isOwnProfile && <>This is your profile</>}
              {profileFollowingLoggedInUser && <>Follows you</>}
            </span>
          )}
        </div>

        <div className="joined">
          <FontAwesomeIcon icon="calendar" />
          Joined {joinedDate.toLocaleString("default", { month: "long" })}{" "}
          {joinedDate.getFullYear()}
        </div>

        <div className="follow-counts">
          <span>
            {followingCount ? (
              <>
                <strong>{followingCount}</strong> Following
              </>
            ) : (
              <>Not following anyone</>
            )}
          </span>
          <span>
            {followersCount ? (
              <>
                <strong>{followersCount}</strong> Followers
              </>
            ) : (
              <>No followers</>
            )}
          </span>
        </div>
      </div>

      <div className="profile-actions">
        {!isOwnProfile && isLoggedIn && (
          <button className="btn btn-light" onClick={handleChangeFollow}>
            {loggedInUserFollowingProfile ? "Unfollow" : "Follow"}
          </button>
        )}
      </div>
    </div>
  );
}
