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

	ReactionCount map[string]int `json:"reactionCount"` // A reaction & count map //
	RetweetCount  int            `json:"retweetCount"`

	UserReactions []string `json:"userReactions"` // A list of reactions //
	UserRetweeted bool     `json:"userRetweeted"`
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
func GetTimeline(userId int, count int) ([]TimelineTweet, map[int]User, error) {
	// N.B First query is to get the tweets, the second one is to get any users that are linked
	// to those tweets (user retweeted or tweeted).

	tweets := make([]TimelineTweet, 0)
	users := make(map[int]User)

	rows, err := store.DB.
		Query(`
			select
				t.id, t.text, t.mediaPath, t.createdAt, 
				u.id, followedTweets.retweeterId
			from tweets t
			join (
				select t.id as tweetId, NULL as retweeterId
					from tweets t
					join follows f on f.followedUserId = t.userId
					join users u on f.followedUserId = u.id
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

	userIds := make([]interface{}, 0)

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

	// [Scaling]
	// Bad to make multiple queries like this.
	// This should be just one query that gets values we're after
	// in one swoop instead of n number of queries.
	// But in this specific use case, it'll do.
	for i := 0; i < len(tweets); i++ {
		retweetCount, err := tweets[i].Tweet.RetweetCount()
		if err != nil {
			return nil, nil, err
		}

		reactionCount, err := tweets[i].Tweet.ReactionCount()
		if err != nil {
			return nil, nil, err
		}

		userReactions, err := tweets[i].Tweet.UserReactions(userId)
		if err != nil {
			return nil, nil, err
		}

		userRetweeted, err := tweets[i].Tweet.UserRetweeted(userId)
		if err != nil {
			return nil, nil, err
		}

		tweets[i].UserReactions = userReactions
		tweets[i].UserRetweeted = userRetweeted
		tweets[i].RetweetCount = retweetCount
		tweets[i].ReactionCount = reactionCount
	}

	if len(userIds) == 0 {
		return tweets, users, nil
	}

	userRows, err := store.DB.
		Query(`
			select u.id, u.username, u.createdAt
			from users u
			where u.id in (?`+strings.Repeat(",?", len(userIds)-1)+`)`,
			userIds...,
		)

	if err != nil {
		return nil, nil, err
	}

	defer userRows.Close()
	for userRows.Next() {
		var u User

		err = userRows.Scan(
			&u.Id, &u.Username, &u.CreatedAt,
		)

		if err != nil {
			return nil, nil, err
		}

		users[u.Id] = u
	}

	return tweets, users, nil
}
