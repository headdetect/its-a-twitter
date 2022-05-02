import React from "react";

import * as AuthContainer from "containers/AuthContainer";
import * as TimeUtils from "utils/TimeUtils";

import "./Tweet.css";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import UserAvatar from "components/UserAvatar";

const API_URL = process.env.REACT_APP_API_URL;

const REACTION_MAP = {
  clap: "👏",
  party: "🎉",
  sad: "😔",
  heart: "❤️",
  thumbsup: "👍",
  thumbsdown: "👎",
};

const ALLOWED_REACTIONS = Object.keys(REACTION_MAP);

export default function Tweet({
  timelineTweet,
  user,
  retweetUser = null,

  onRemoveRetweet = _ => {},
  onRetweet = _ => {},
  onRemoveReaction = (_, __) => {},
  onReaction = (_, __) => {},
  onDeleteTweet = _ => {},
}) {
  const { loggedInUser, isLoggedIn } = AuthContainer.useContext();
  const [titles, setTitles] = React.useState({
    reaction: "Add a reaction",
    share: "Share this tweet",
    retweet: "Retweet",
  });

  const isOwnTweet = loggedInUser && loggedInUser.id === user.id;

  React.useEffect(() => {
    if (!loggedInUser) {
      setTitles(t => ({
        ...t,
        reaction: "Must be logged in to react",
        retweet: "Must be logged in to retweet",
      }));

      return;
    }

    if (isOwnTweet) {
      setTitles(t => ({
        ...t,
        reaction: "Cannot react to your own tweet",
        retweet: "Cannot retweet your own tweet",
      }));
    }
  }, [loggedInUser, isOwnTweet]);

  const {
    tweet,
    reactionCount = {},
    retweetCount = 0,
    userReactions = [],
    userRetweeted = false,
  } = timelineTweet;

  const handleRetweet = async () => {
    if (userRetweeted) {
      onRemoveRetweet(tweet.id);
    } else {
      onRetweet(tweet.id);
    }
  };

  const handleReact = r => {
    if (userReactions.includes(r)) {
      onRemoveReaction(tweet.id, r);
    } else {
      onReaction(tweet.id, r);
    }
  };

  const handleDeleteTweet = () => {
    onDeleteTweet(tweet.id);
  };

  return (
    <>
      {retweetUser && (
        <div className="retweet">
          <UserAvatar
            user={retweetUser}
            size={25}
            style={{ marginRight: "var(--spacing)" }}
          />

          <a
            href={`/profile/@${retweetUser.username}`}
            rel="noopener noreferrer"
          >
            @{retweetUser.username}
          </a>
          <span>retweeted</span>
        </div>
      )}
      <div className="tweet">
        {isOwnTweet && (
          <button
            className="btn btn-delete-tweet"
            title="Delete tweet"
            onClick={handleDeleteTweet}
          >
            <FontAwesomeIcon icon="trash" />
          </button>
        )}

        <UserAvatar
          user={user}
          style={{ marginRight: "calc(var(--spacing) * 2.5)" }}
        />
        <div className="tweet-content">
          <div className="user-info">
            <a className="user-link" href={`/profile/@${user.username}`}>
              @{user.username}
            </a>
            <a className="tweet-link" href={`/tweet/${tweet.id}`}>
              {TimeUtils.toAgoString(new Date(tweet.createdAt * 1000))}
            </a>
          </div>

          <p>{tweet.text}</p>

          {tweet.mediaPath && (
            <img
              src={`${API_URL}/asset/${tweet.mediaPath}`}
              alt="tweet media"
            />
          )}

          <div className="tweet-actions">
            <button
              onClick={handleRetweet}
              className={`btn btn-tweet-action ${
                userRetweeted ? "selected" : ""
              }`}
              disabled={!isLoggedIn || isOwnTweet}
              title={titles.retweet}
            >
              <FontAwesomeIcon icon="retweet" /> <span>{retweetCount}</span>
            </button>

            <div className="tweet-reactions">
              {ALLOWED_REACTIONS.map(r => (
                <button
                  key={r}
                  disabled={!isLoggedIn || isOwnTweet}
                  onClick={() => handleReact(r)}
                  className={`btn btn-tweet-reaction ${
                    userReactions.includes(r) ? "selected" : ""
                  }`}
                  title={titles.reaction}
                >
                  {REACTION_MAP[r]} {reactionCount[r] || 0}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
