import React from "react";
import { useHref } from "react-router-dom";

import { API_URL } from "consts";

import "./UserAvatar.css";

export default function UserAvatar({
  user,

  size = 50,
  imageProps = { style: {} },
  onClick = undefined,
  ...theRest
}) {
  const defaultHrefPath = useHref(`/profile/@${user.username}`);

  const hrefOrClick = {};

  if (onClick) {
    hrefOrClick.onClick = onClick;
  } else {
    hrefOrClick.href = defaultHrefPath;
  }

  const imageSrc = user.profilePicPath.startsWith("data:image")
    ? user.profilePicPath // Use the raw base64 image
    : `${API_URL}/asset/${user.profilePicPath}`; // Treat as a url path

  return (
    <a
      {...hrefOrClick}
      {...theRest}
      style={{
        width: size,
        height: size,
        ...theRest.style,
      }}
    >
      <img
        className="avatar"
        src={imageSrc}
        alt={`${user.username}'s profile`}
        {...imageProps}
        style={{
          width: size,
          height: size,
          ...imageProps?.style,
        }}
      />
    </a>
  );
}
