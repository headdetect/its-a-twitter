/**
 * @fileoverview Container that gives the ability to tweet, retweet, post 
 * media, fetch timelines. Do anything surrounding tweeting.
 */
import React from "react";
import PropTypes from "prop-types";
import useHttpService from "../hooks/useHttpService";
import * as UserContainer from "./UserContainer";

export const Context = React.createContext(null);

export function Provider(props) {
  const { authenticatedFetch } = useHttpService();
  const userContext = UserContainer.useContext();

  const tweet = React.useCallback((text, media) => {

  }, []);

  const deleteTweet = React.useCallback((tweetId) => {

  }, []);

  const retweet = React.useCallback((tweetId) => {

  }, []);

  const fetchTimeline = React.useCallback(() => {

  }, []);

  const reactToTweet = React.useCallback((tweetId) => {

  }, []);

  return (
    <Context.Provider
      value={{
        tweet,
        deleteTweet,
        retweet,
        fetchTimeline,
        reactToTweet
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
      "container can only be used in the context of a TimelineContainer.Provider"
    );
  }

  return context;
}

Provider.propTypes = {
  children: PropTypes.element.isRequired,
};