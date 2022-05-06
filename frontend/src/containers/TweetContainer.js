/**
 * @fileoverview Container that gives the ability to tweet, retweet, post
 * media, fetch timelines. Do anything surrounding tweeting.
 */
import React from "react";
import * as AuthContainer from "containers/AuthContainer";

import { API_URL } from "consts";

export const Context = React.createContext(null);

export function Provider({ children }) {
  const { authenticatedFetch, loggedInUser, isLoggedIn } =
    AuthContainer.useContext();

  const [timeline, setTimeline] = React.useState([]);
  const [timelineUsers, setTimelineUsers] = React.useState([]);
  const [timelineStatus, setTimelineStatus] = React.useState("loading"); // loading | finished | error | not-logged-in

  const fetchAndStoreTimeline = React.useCallback(async () => {
    // We're okay with this running multiple times. It should reload
    // the timeline every time the login state is set to true
    try {
      let response;
      if (isLoggedIn)
        response = await authenticatedFetch(`${API_URL}/timeline`);
      else response = await fetch(`${API_URL}/timeline`);

      if (response.status === 401) {
        setTimeline([]);
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
  }, [authenticatedFetch, isLoggedIn]);

  const submitTweet = React.useCallback(
    async (text, media) => {
      const formData = new FormData();
      formData.append("file", media);
      formData.append("text", text);

      const response = await authenticatedFetch(`${API_URL}/tweet`, {
        method: "POST",
        body: formData,
      });

      if (response.status === 401) {
        // TODO: Make the user log-in
        return;
      }

      const body = await response.json();

      if (!body.tweet) {
        throw new Error("Server returned something unexpected");
      }

      const timelineTweet = {
        // Defaults //
        posterUserId: loggedInUser.id,
        reactionCount: {},
        retweetCount: 0,
        retweeterUserId: null,
        userReactions: [],
        userRetweeted: false,

        // Results //
        ...body.tweet,
      };

      setTimelineUsers(oldTimelineUsers => ({
        ...oldTimelineUsers,
        [`${loggedInUser.id}`]: loggedInUser,
      }));

      setTimeline(oldTimeline => [timelineTweet, ...oldTimeline]);
    },
    [authenticatedFetch, loggedInUser],
  );

  const getTweet = React.useCallback(
    async tweetId => {
      let response;
      if (isLoggedIn)
        response = await authenticatedFetch(`${API_URL}/tweet/${tweetId}`);
      else response = await fetch(`${API_URL}/tweet/${tweetId}`);

      const body = await response.json();

      if (body.tweet) {
        return body;
      }

      throw new Error("Malformed response was returned from the server");
    },
    [authenticatedFetch, isLoggedIn],
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

      setTimeline(oldTimeline =>
        oldTimeline.filter(t => t.tweet.id !== tweetId),
      );
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

      setTimeline(line =>
        line.map(t => {
          if (t.tweet.id !== tweetId) return t;

          return {
            ...t,
            userRetweeted: true,
            retweetCount: (t.retweetCount || 0) + 1,
          };
        }),
      );
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

      setTimeline(line =>
        line.map(t => {
          if (t.tweet.id !== tweetId) return t;

          return {
            ...t,
            userRetweeted: false,
            retweetCount: t.retweetCount - 1,
          };
        }),
      );
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

      setTimeline(line =>
        line.map(t => {
          if (t.tweet.id !== tweetId) return t;

          return {
            ...t,
            userReactions: [
              ...(t.userReactions?.filter(r => r !== reaction) ?? []), // filter out any in case they exists //
              reaction,
            ],
            reactionCount: {
              ...t.reactionCount,
              [`${reaction}`]: (t.reactionCount?.[reaction] || 0) + 1,
            },
          };
        }),
      );
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

      setTimeline(line =>
        line.map(t => {
          if (t.tweet.id !== tweetId) return t;

          return {
            ...t,
            userReactions: t.userReactions?.filter(r => r !== reaction) ?? [],
            reactionCount: {
              ...t.reactionCount,
              [`${reaction}`]: t.reactionCount[reaction] - 1,
            },
          };
        }),
      );
    },
    [authenticatedFetch],
  );

  return (
    <Context.Provider
      value={{
        // Actions //
        submitTweet,
        getTweet,
        deleteTweet,
        retweet,
        removeRetweet,
        fetchAndStoreTimeline,
        addReaction,
        removeReaction,
        setTimeline,

        // States //
        timeline,
        timelineUsers,
        timelineStatus,
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
      "container can only be used in the context of a TweetContainer.Provider",
    );
  }

  return context;
}
