package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/headdetect/its-a-twitter/api/store"
)

var PROFILE_PIC_PATH = "./assets/profile"

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	CreatedAt int64 `json:"createdAt"`
}

type Follow struct {
	User         *User `json:"user"`
	FollowedUser *User `json:"followedUser"`

	CreatedAt int64 `json:"createdAt"`
}

func GetUserByUsernameWithPass(username string) (user User, hashedPassword string, email string, err error) {
	err = store.DB.QueryRow(
		"select id, username, email, password, createdAt from users where username = ? limit 1",
		username,
	).Scan(
		&user.Id, &user.Username, &email, &hashedPassword, &user.CreatedAt,
	)

	return
}

func GetUserById(id int) (user User, email string, err error) {
	err = store.DB.QueryRow(
		"select id, username, email, createdAt from users where id = ? limit 1",
		id,
	).Scan(
		&user.Id, &user.Username, &email, &user.CreatedAt,
	)

	return
}

func GetUserByUsername(username string) (user User, email string, err error) {
	err = store.DB.QueryRow(
		"select id, username, email, createdAt from users where username = ? limit 1",
		username,
	).Scan(
		&user.Id, &user.Username, &user.Email, &user.CreatedAt,
	)

	return
}

func GetUserByTweetId(tweetId int) (user User, email string, err error) {
	err = store.DB.QueryRow(
		"select u.id, u.username, u.email, u.createdAt from users u, tweets t where u.id = t.userId and t.id = 1 limit 1",
	).Scan(
		&user.Id, &user.Username, &user.Email, &user.CreatedAt,
	)

	return
}

func MakeUser(email string, username string, passwordHash string) (user User, err error) {
	createdAt := time.Now().Unix()

	res, err := store.DB.Exec(
		"insert into users (username, password, email, createdAt) values (?, ?, ?, ?)",
		username, passwordHash, email, createdAt,
	)

	if err != nil {
		return
	}

	id, err := res.LastInsertId()

	if err != nil {
		return
	}

	user.Id = int(id)
	user.CreatedAt = createdAt
	user.Username = username

	return
}

func (u *User) DeleteUser() (err error) {
	_, err = store.DB.Exec(
		"delete from users where id = ?",
		u.Id,
	)

	return
}

func (u *User) GetProfilePicPath() (string, error) {
	// TODO: Look for file

	return fmt.Sprintf("%s/%d.jpg", PROFILE_PIC_PATH, u.Id), nil
}

func (u *User) FollowUser(followedUserId int) (err error) {
	_, err = store.DB.Exec(
		"insert into follows (userId, followedUserId) values (?, ?)",
		u.Id, followedUserId,
	)

	return
}

func (u *User) UnFollowUser(followedUserId int) (err error) {
	_, err = store.DB.Exec(
		"delete from follows where userId = ? and followedUserId = ?",
		u.Id, followedUserId,
	)

	return
}

func (u *User) GetFollowers() ([]User, error) {
	rows, err := store.DB.
		Query(
			`select 
				followers.id,
				followers.username,
				followers.createdAt
			from follows f
			join users followers on f.userId = followers.id
			where f.followedUserId = ?`,
			u.Id,
		)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	followers := []User{}

	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Username, &u.CreatedAt)

		if err != nil {
			return nil, err
		}

		followers = append(followers, u)
	}

	return followers, nil
}

func (u *User) GetFollowing() ([]User, error) {
	rows, err := store.DB.
		Query(
			`select 
				following.id,
				following.username,
				following.createdAt
			from follows
			join users following on follows.followedUserId = following.id
			where follows.userId = ?`,
			u.Id,
		)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	following := []User{}

	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Username, &u.CreatedAt)

		if err != nil {
			return nil, err
		}

		following = append(following, u)
	}

	return following, nil
}

func (u *User) GetTweets() ([]Tweet, error) {
	rows, err := store.DB.
		Query(`select id, text, mediaPath, createdAt from tweets where userId = ?`,
			u.Id,
		)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tweets := []Tweet{}

	for rows.Next() {
		var t Tweet
		var mediaPath sql.NullString
		err := rows.Scan(&t.Id, &t.Text, &mediaPath, &t.CreatedAt)

		if err != nil {
			return nil, err
		}

		if mediaPath.Valid {
			t.MediaPath = mediaPath.String
		}

		// Note: We don't attach the user to the tweet, because the caller
		// should already have access to the user

		tweets = append(tweets, t)
	}

	return tweets, nil
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
func (u *User) GetTimeline(count int) ([]TimelineTweet, []User, error) {
	// N.B First query is to get the tweets, the second one is to get any users that are linked
	// to those tweets (user retweeted or tweeted).

	tweets := make([]TimelineTweet, 0)
	users := make([]User, 0)

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
					where f.userId = $1
				union
				select rt.tweetId as tweetId, f.followedUserId as retweeterId
					from follows f
					join retweets rt on rt.userId = f.followedUserId
					where f.userId = $1
			) followedTweets on followedTweets.tweetId = t.id
			join users u on u.id = t.userId
			where t.userId != $1
			group by t.id
			order by t.createdAt desc`,
			u.Id,
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
	for _, t := range tweets {
		retweetCount, err := t.Tweet.RetweetCount()
		if err != nil {
			return nil, nil, err
		}

		reactionCount, err := t.Tweet.ReactionCount()
		if err != nil {
			return nil, nil, err
		}

		t.RetweetCount = retweetCount
		t.ReactionCount = reactionCount
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

		users = append(users, u)
	}

	return tweets, users, nil
}
