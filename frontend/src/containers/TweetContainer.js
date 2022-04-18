/**
 * @fileoverview Container that gives the ability to tweet, retweet, post
 * media, fetch timelines. Do anything surrounding tweeting.
 */
import React from "react";
import * as AuthContainer from "containers/AuthContainer";
import * as UserContainer from "containers/UserContainer";

const API_URL = process.env.REACT_APP_API_URL;

export const Context = React.createContext(null);

export function Provider(props) {
  const { authenticatedFetch } = AuthContainer.useContext();
  const { isLoggedIn } = UserContainer.useContext();

  const [timeline, setTimeline] = React.useState(undefined);
  const [timelineUsers, setTimelineUsers] = React.useState(undefined);
  const [timelineStatus, setTimelineStatus] = React.useState("init"); // init | loading | finished | error

  const refreshTimeline = React.useCallback(
    async filter => {
      setTimelineStatus("loading");

      try {
        const response = await authenticatedFetch(
          `${API_URL}/timeline${filter ? `?filter=${filter}` : ""}`,
        );

        if (response.status === 401) {
          setTimelineStatus("error");
          return;
        }

        const obj = await response.json();

        if (obj.tweets && obj.users) {
          setTimeline(obj.tweets);
          setTimelineUsers(obj.users);

          setTimelineStatus("finished");
        } else {
          setTimelineStatus("error");
        }
      } catch (e) {
        setTimelineStatus("error");
      }
    },
    [authenticatedFetch],
  );

  const tweet = React.useCallback(
    async (text, media) => {
      // TODO: Post medias

      const response = await authenticatedFetch(`${API_URL}/tweet`, {
        method: "POST",
        body: JSON.stringify({ text }),
      });

      if (response.status === 401) {
        // TODO: Make the user log-in
        return;
      }
    },
    [authenticatedFetch],
  );

  const deleteTweet = React.useCallback(
    async tweetId => {
      const response = await authenticatedFetch(`${API_URL}/tweet/${tweetId}`, {
        method: "DELETE",
      });

      if (response.status === 401) {
        // TODO: Make the user log-in
        return;
      }

      if (!response.ok) {
        throw new Error("Error deleting tweet");
      }
    },
    [authenticatedFetch],
  );

  const retweet = React.useCallback(
    async tweetId => {
      const response = await authenticatedFetch(
        `${API_URL}/tweet/${tweetId}/retweet`,
        {
          method: "PUT",
        },
      );

      if (response.status === 401) {
        // TODO: Make the user log-in
        return;
      }

      if (!response.ok) {
        throw new Error("Error retweeting");
      }
    },
    [authenticatedFetch],
  );

  const removeRetweet = React.useCallback(
    async tweetId => {
      const response = await authenticatedFetch(
        `${API_URL}/tweet/${tweetId}/retweet`,
        {
          method: "DELETE",
        },
      );

      if (response.status === 401) {
        // TODO: Make the user log-in
        return;
      }

      if (!response.ok) {
        throw new Error("Error removing retweet");
      }
    },
    [authenticatedFetch],
  );

  const addReaction = React.useCallback(
    async (tweetId, reaction) => {
      const response = await authenticatedFetch(
        `${API_URL}/tweet/${tweetId}/react/${reaction}`,
        {
          method: "PUT",
        },
      );

      if (response.status === 401) {
        // TODO: Make the user log-in
        return;
      }

      if (!response.ok) {
        throw new Error("Error adding reaction");
      }
    },
    [authenticatedFetch],
  );

  const removeReaction = React.useCallback(
    async (tweetId, reaction) => {
      const response = await authenticatedFetch(
        `${API_URL}/tweet/${tweetId}/react/${reaction}`,
        {
          method: "DELETE",
        },
      );

      if (response.status === 401) {
        // TODO: Make the user log-in
        return;
      }

      if (!response.ok) {
        throw new Error("Error removing reaction");
      }
    },
    [authenticatedFetch],
  );

  React.useEffect(() => {
    // We're okay with this running multiple times. It should reload
    // the timeline every time the login state is set to true

    if (!isLoggedIn) {
      // TODO: Remove this when we have un-auth'd timelines
      return;
    }

    refreshTimeline();
  }, [refreshTimeline, isLoggedIn]);

  return (
    <Context.Provider
      value={{
        // Actions //
        tweet,
        deleteTweet,
        retweet,
        removeRetweet,
        refreshTimeline,
        addReaction,
        removeReaction,

        // States //
        timeline,
        timelineUsers,
        timelineStatus,
      }}
    >
      {props.children}
    </Context.Provider>
  );
}

export function useContext() {
  const context = React.useContext(Context);

  if (!context) {
    throw new Error(
      "container can only be used in the context of a TweetContainer.Provider",
    );
  }

  return context;
}
