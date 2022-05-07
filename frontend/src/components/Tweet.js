import React from "react";
import { Link } from "react-router-dom";

import * as AuthContainer from "containers/AuthContainer";
import * as TimeUtils from "utils/TimeUtils";
import * as UrlUtils from "utils/UrlUtils";
import * as StringUtils from "utils/StringUtils";

import "./Tweet.css";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import UserAvatar from "components/UserAvatar";

const API_URL = process.env.REACT_APP_API_URL;

const REACTION_MAP = {
  laugh: "😂",
  party: "🎉",
  sad: "😔",
  heart: "❤️",
  raisedHands: "🙌",
  shocked: "😲",
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

  const linkedText = UrlUtils.insertAnchorElements(tweet.text);
  const newlinedText = StringUtils.insertNewlineElements(linkedText);

  return (
    <>
      {retweetUser && (
        <div className="retweet">
          <UserAvatar
            user={retweetUser}
            size={25}
            style={{ marginRight: "var(--spacing)" }}
          />

          <Link
            to={`/profile/@${retweetUser.username}`}
            rel="noopener noreferrer"
          >
            @{retweetUser.username}
          </Link>
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

        <div className="tweet-info">
          <UserAvatar
            user={user}
            style={{ marginRight: "calc(var(--spacing) * 2.5)" }}
          />

          <div className="user-info">
            <Link className="user-link" to={`/profile/@${user.username}`}>
              @{user.username}
            </Link>
            <Link className="tweet-link" to={`/tweet/${tweet.id}`}>
              {TimeUtils.toAgoString(new Date(tweet.createdAt * 1000))}
            </Link>
          </div>
        </div>

        <div className="tweet-content">
          <p>{newlinedText}</p>

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
