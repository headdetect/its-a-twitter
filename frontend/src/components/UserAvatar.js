import React from "react";

import "./UserAvatar.css";

export default function UserAvatar({
  user,
  size = 50,
  imageProps = { style: {} },
  ...theRest
}) {
  //TODO: Lazy load the images in as they're needed

  return (
    <a
      href={`/profile/@${user.username}`}
      {...theRest}
      style={{
        width: size,
        height: size,
        ...theRest.style,
      }}
    >
      <img
        className="avatar"
        src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAJYAAACWCAAAAAAZai4+AAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAAmJLR0QA/4ePzL8AAAAHdElNRQfmBBwKOCIs7mQRAAACaElEQVR42u3V3XKiMAAF4L7/G5xErYBaq9v6s6OtLlD/ocBDbQCB7Ey7FbgwF+fcSIA5fgYTHhIj83BvAFlkkUUWWWQZErLIIsuEkEUWWSaELLLIMiFkkUWWCSGLLLJMCFlkkWVCyCKLLBNCFllkmRCyyLoD6wXXgz2K+Okwfht05OAtvqm/Ycn3rM9O0bj9pzEa5cej6AZV05JvWeEYReMSIz9PqEYzyNUpeJOY/axqXPI1K5g/SZSNUyy0S4Cbfv6BCP5valPyNWuXT/F1NMCmurRCPz94xKqclB5+ZwcvsOOmJTewouPx6JaNHeyrS8+Y5wdzTMqTLsRRfXzkHw1LfmalORWNIfDx3BP2r3M6srDNT29hVTdPMYiTqH+dtIYl9Vjl0pbv2a/289M+utXNYRfrZJHaWpTUY20A2wsCz84ekcAhP32A0O52IV0hTu1KarH812W2u8Q2xkkiq0ap3z5VM7FqW1KHVcZLZ7xfzX9fv6geoxO3LWnEUrtNkAyrf+tQv3gQkJe2JY1YYdo4KbbFJabatchCF08tS+qwYutxW8y4mv91saQttfaqLOCcpb5jNimpNVtjWJ9ZtYPXJLkI7NLRDkJ7Zvt0ga3RCdqU1GOpLcfxgovnoJt2zNDzosTv6W9Z9QgX2VdOWpTUZCUbke+E3ewXRkNAqHfwSFt3c1jp8ld/e7d5SV1Wcp5ZUtqLMB/Fa0dKZ60Vqpfy7srrhU1LbmTdNWSRRZYJIYssskwIWWSRZULIIossE0IWWWSZELLIIsuEkEUWWSaELLLIMiFkkUWWCSGLLLJMCFlkkWVC/gI+i6YMw2buTAAAACV0RVh0ZGF0ZTpjcmVhdGUAMjAyMi0wNC0yOFQxOTo1NjozNCswOTowMA2Q9vkAAAAldEVYdGRhdGU6bW9kaWZ5ADIwMjItMDQtMjhUMTk6NTY6MzQrMDk6MDB8zU5FAAAAAElFTkSuQmCC"
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
