package model

import (
	"database/sql"
	"strings"

	"github.com/headdetect/its-a-twitter/api/store"
)

type TimelineTweet struct {
	Tweet Tweet `json:"tweet"`

	Poster    int  `json:"posterUserId"`
	Retweeter *int `json:"retweeterUserId"`

	ReactionCount map[string]int `json:"reactionCount,omitempty"` // A reaction & count map //
	RetweetCount  int            `json:"retweetCount,omitempty"`

	UserReactions []string `json:"userReactions,omitempty"` // A list of reactions //
	UserRetweeted bool     `json:"userRetweeted,omitempty"`
}

func fillTimelineMeta(tweets *[]TimelineTweet, asUserId int, showUserInteractions bool) (err error) {
	pTweets := *tweets

	// [Scaling]
	// Bad to make multiple queries like this.
	// This should be just one query that gets values we're after
	// in one swoop instead of n number of queries.
	// But in this specific use case, it'll do.
	for i := 0; i < len(pTweets); i++ {
		retweetCount, err := (pTweets)[i].Tweet.RetweetCount()
		if err != nil {
			return err
		}

		reactionCount, err := pTweets[i].Tweet.ReactionCount()
		if err != nil {
			return err
		}

		if showUserInteractions {
			userReactions, err := pTweets[i].Tweet.UserReactions(asUserId)
			if err != nil {
				return err
			}

			userRetweeted, err := pTweets[i].Tweet.UserRetweeted(asUserId)
			if err != nil {
				return err
			}

			pTweets[i].UserReactions = userReactions
			pTweets[i].UserRetweeted = userRetweeted
		} else {
			pTweets[i].UserReactions = make([]string, 0)
			pTweets[i].UserRetweeted = false
		}

		pTweets[i].RetweetCount = retweetCount
		pTweets[i].ReactionCount = reactionCount
	}

	return nil
}

func fillAssociatedUsers(users *map[int]User, userIds []any) (err error) {
	if len(userIds) == 0 {
		return nil
	}

	userRows, err := store.DB.
		Query(`
			select u.id, u.username, u.createdAt
			from users u
			where u.id in (?`+strings.Repeat(",?", len(userIds)-1)+`)`,
			userIds...,
		)

	if err != nil {
		return err
	}

	defer userRows.Close()
	for userRows.Next() {
		var u User

		if err = userRows.Scan(
			&u.Id, &u.Username, &u.CreatedAt,
		); err != nil {
			return err
		}

		(*users)[u.Id] = u
	}

	return nil
}

// A generated timeline for those that are not logged in. Will just sort by most recent
func GetFeatured() ([]TimelineTweet, map[int]User, error) {
	// N.B First query is to get the tweets, the second one is to get any users that are linked
	// to those tweets (user retweeted or tweeted).
	tweets := make([]TimelineTweet, 0)
	users := make(map[int]User)

	rows, err := store.DB.
		Query(`
			select
				t.id, t.text, t.mediaPath, t.createdAt, 
				u.id, u.username, u.createdAt
			from tweets t
			join users u on t.userId = u.id
			order by t.createdAt DESC`,
		)

	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t TimelineTweet
		var u User
		var mediaPath sql.NullString

		err = rows.Scan(
			&t.Tweet.Id, &t.Tweet.Text, &mediaPath, &t.Tweet.CreatedAt, // Tweet properties //
			&u.Id, &u.Username, &u.CreatedAt, // User properties //
		)

		if err != nil {
			return nil, nil, err
		}

		if mediaPath.Valid {
			t.Tweet.MediaPath = mediaPath.String
		}

		t.Poster = u.Id
		t.Tweet.User = u

		// N.B: We don't attach the user to the tweet, because the caller
		// should already have access to the user via the users list
		tweets = append(tweets, t)
		users[u.Id] = u
	}

	if err = fillTimelineMeta(&tweets, 0, false); err != nil {
		return nil, nil, err
	}

	return tweets, users, nil
}

// Get user's stream of tweets we want to feature
// The timeline should consist of tweets from followed
// users and from retweets of followed users.
//
// The for retweets to show up, the user does not need to be
// following the original poster of the retweeted tweet.
//
// Should only show the `count` latest tweets.
// Should be in order of `createdAt`
func GetTimeline(userId int) ([]TimelineTweet, map[int]User, error) {
	// N.B First query is to get the tweets, the second one is to get any users that are linked
	// to those tweets (user retweeted or tweeted).
	tweets := make([]TimelineTweet, 0)
	users := make(map[int]User)
	userIds := make([]any, 0) // The type is any, but they're really ints //

	rows, err := store.DB.
		Query(`
			select
				t.id, t.text, t.mediaPath, t.createdAt, 
				u.id, followedTweets.retweeterId
			from tweets t
			join (
				select t.id as tweetId, NULL as retweeterId
					from tweets t
					left join follows f on f.followedUserId = t.userId
					left join users u on f.followedUserId = u.id
					where f.userId = $1 or t.userId = $1
				union
				select rt.tweetId as tweetId, f.followedUserId as retweeterId
					from follows f
					join retweets rt on rt.userId = f.followedUserId
					where f.userId = $1
			) followedTweets on followedTweets.tweetId = t.id or t.userId = $1
			join users u on u.id = t.userId
			group by t.id
			order by t.createdAt desc`,
			userId,
		)

	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t TimelineTweet
		var mediaPath sql.NullString
		var retweeterId sql.NullInt64
		var tweeterId int

		err = rows.Scan(
			&t.Tweet.Id, &t.Tweet.Text, &mediaPath, &t.Tweet.CreatedAt, // Tweet properties //
			&tweeterId, &retweeterId, // User ids //
		)

		if err != nil {
			return nil, nil, err
		}

		if mediaPath.Valid {
			t.Tweet.MediaPath = mediaPath.String
		}

		if retweeterId.Valid {
			// [Scaling]
			// TODO
			rId := int(retweeterId.Int64)
			userIds = append(userIds, rId)
			t.Retweeter = &rId
		}

		t.Poster = tweeterId

		// N.B: We don't attach the user to the tweet, because the caller
		// should already have access to the user via the users list
		tweets = append(tweets, t)
		userIds = append(userIds, tweeterId)
	}

	if err = fillTimelineMeta(&tweets, userId, true); err != nil {
		return nil, nil, err
	}

	if err = fillAssociatedUsers(&users, userIds); err != nil {
		return nil, nil, err
	}

	return tweets, users, nil
}

func GetTimelineByUser(user User, loggedInUser *User) ([]TimelineTweet, map[int]User, error) {
	// N.B First query is to get the tweets, the second one is to get any users that are linked
	// to those tweets (user retweeted or tweeted).
	tweets := make([]TimelineTweet, 0)
	users := make(map[int]User)
	userIds := make([]interface{}, 0)

	/*
	 * Get all the tweets and retweets from user
	 */
	rows, err := store.DB.
		Query(`
			select 
				t.id, t.text, t.mediaPath, t.createdAt, 
				t.userId
			from tweets t
			left join retweets r on r.tweetId = t.id
			where t.userId = $1 or r.userId = $1
			group by t.id
			order by t.createdAt DESC`,
			user.Id,
		)

	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t TimelineTweet
		var mediaPath sql.NullString

		err = rows.Scan(
			&t.Tweet.Id, &t.Tweet.Text, &mediaPath, &t.Tweet.CreatedAt, // Tweet properties //
			&t.Poster,
		)

		if err != nil {
			return nil, nil, err
		}

		if mediaPath.Valid {
			t.Tweet.MediaPath = mediaPath.String
		}

		if t.Poster != user.Id {
			// This is a retweet //
			userIds = append(userIds, t.Poster)
		}

		// N.B: We don't attach the user to the tweet, because the caller
		// should already have access to the user via the users list
		tweets = append(tweets, t)
	}

	if loggedInUser != nil {
		err = fillTimelineMeta(&tweets, loggedInUser.Id, true)
	} else {
		err = fillTimelineMeta(&tweets, 0, false)
	}

	if err != nil {
		return nil, nil, err
	}

	if err = fillAssociatedUsers(&users, userIds); err != nil {
		return nil, nil, err
	}

	return tweets, users, nil
}
