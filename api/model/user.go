package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/headdetect/its-a-twitter/api/store"
)

var PROFILE_PIC_PATH = "./assets/profile"

type User struct {
	Id int
	Username string

	CreatedAt int64
}

type Follow struct {
	Id int

	User *User
	FollowedUser *User

	CreatedAt int64
}

func GetUserWithPassByUsername(username string) (User, string, error) {
	var user User
	var hashedPassword string

	err := store.DB.QueryRow(
			"select id, username, password, createdAt from users where username = ? limit 1",
			username,
		).Scan(
			&user.Id, &user.Username, &hashedPassword, &user.CreatedAt,
		)

	return user, hashedPassword, err
}

func GetUserById(id int) (User, error) {
	var user User

	err := store.DB.QueryRow(
			"select id, username, createdAt from user where id = ? limit 1", 
			id,
		).Scan(
			&user.Id, &user.Username, &user.CreatedAt,
		)

	return user, err
}

func GetUserByUsername(username string) (User, error) {
	var user User

	err := store.DB.QueryRow(
			"select id, username, createdAt from users where username = ? limit 1", 
			username,
		).Scan(
			&user.Id, &user.Username, &user.CreatedAt,
		)

	return user, err
}

func MakeUser(username string, passwordHash string) (User, error) {
	var user User
	
	createdAt := time.Now().Unix()

	res, err := store.DB.Exec(
		"insert into users (username, password, createdAt) values (?, ?, ?, ?)",
		username, passwordHash, createdAt,
	)

	if err != nil {
		return user, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return user, err
	}

	user.Id = int(id)
	user.CreatedAt = createdAt
	user.Username = username

	return user, err
}

func (u *User) DeleteUser() (error) {
	_, err := store.DB.Exec(
		"delete from users where id = ?",
		u.Id,
	)

	return err
}

func (u *User) GetProfilePicPath() (string, error) {
	// TODO: Look for file

	return fmt.Sprintf("%s/%d.jpg", PROFILE_PIC_PATH, u.Id), nil
}


func (u *User) FollowUser(followedUserId int) (error) {
	_, err := store.DB.Exec(
		"insert into follows (userId, followedUserId) values (?, ?)",
		u.Id, followedUserId,
	)

	return err
}

func (u *User) UnFollowUser(followedUserId int) (error) {
	_, err := store.DB.Exec(
		"delete from follows where userId = ? and followedUserId = ?",
		u.Id, followedUserId,
	)

	return err
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

	followering := []User{}
	
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Username, &u.CreatedAt)

		if err != nil {
			return nil, err
		}

		followering = append(followering, u)
	}

	return followering, nil
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

		if (mediaPath.Valid) {
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
func (u *User) GetTimeline(count int) ([]Tweet, error) {
	rows, err := store.DB.
		Query(`
			SELECT t.id, t.text, t.mediaPath, t.createdAt
			FROM tweets t
			JOIN (
				SELECT t.id as tweetId
					FROM tweets t
					JOIN follows f ON f.followedUserId = t.userId
					JOIN users u ON f.followedUserId = u.id
					WHERE f.userId = ?
				UNION
				SELECT rt.tweetId as tweetId
					FROM follows f
					JOIN retweets rt ON rt.userId = f.followedUserId
					WHERE f.userId = ?
			) followedTweets ON followedTweets.tweetId = t.id
			ORDER BY t.createdAt DESC`,

			// The userId is filled in two spots. Instead of 
			// messing with named parameters, I'll just fill twice
			u.Id, u.Id,
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

		if (mediaPath.Valid) {
			t.MediaPath = mediaPath.String
		}
		
		// Note: We don't attach the user to the tweet, because the caller
		// should already have access to the user

		tweets = append(tweets, t)
	}

	return tweets, nil
}