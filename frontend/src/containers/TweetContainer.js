/**
 * @fileoverview Container that gives the ability to tweet, retweet, post
 * media, fetch timelines. Do anything surrounding tweeting.
 */
import React from "react";
import * as AuthContainer from "containers/AuthContainer";

const API_URL = process.env.REACT_APP_API_URL;

export const Context = React.createContext(null);

export function Provider({ children }) {
  const { authenticatedFetch, loggedInUser, isLoggedIn } =
    AuthContainer.useContext();

  const [timeline, setTimeline] = React.useState(undefined);
  const [timelineUsers, setTimelineUsers] = React.useState(undefined);
  const [timelineStatus, setTimelineStatus] = React.useState("loading"); // loading | finished | error | not-logged-in

  const refreshTimeline = React.useCallback(async () => {
    // We're okay with this running multiple times. It should reload
    // the timeline every time the login state is set to true
    try {
      let response;
      if (isLoggedIn)
        response = await authenticatedFetch(`${API_URL}/timeline`);
      else response = await fetch(`${API_URL}/timeline`);

      if (response.status === 401) {
        debugger;
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

      const body = await response.json();

      if (!body.tweet) {
        throw new Error("Server returned something unexpected");
      }

      const timelineTweet = {
        tweet: body.tweet,
        posterUserId: loggedInUser.id,
        reactionCount: {},
        retweetCount: 0,
        retweeterUserId: null,
        userReactions: [],
        userRetweeted: false,
      };

      setTimeline(oldTimeline => [timelineTweet, ...oldTimeline]);
    },
    [authenticatedFetch, loggedInUser],
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
            retweetCount: t.retweetCount + 1,
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
              ...t.userReactions.filter(r => r !== reaction), // filter out any in case they exists //
              reaction,
            ],
            reactionCount: {
              ...t.reactionCount,
              [`${reaction}`]: (t.reactionCount[reaction] || 0) + 1,
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
            userReactions: t.userReactions.filter(r => r !== reaction),
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
        tweet,
        deleteTweet,
        retweet,
        removeRetweet,
        refreshTimeline,
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
