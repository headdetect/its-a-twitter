package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/headdetect/its-a-twitter/api/store"
	"github.com/headdetect/its-a-twitter/api/utils"
)

type Tweet struct {
	Id int `json:"id"`
	User *User `json:"user"`
	Text string `json:"text"`
	MediaPath string `json:"mediaPath"` // TODO: Do we replace this with a hashed version of the id?

	CreatedAt int64 `json:"createdAt"`
}

type Retweet struct {
	Id int `json:"id"`
	Tweet *Tweet `json:"tweet"`
	User *User `json:"user"`

	CreatedAt int64 `json:"createdAt"`
}

type Reaction struct {
	Id int `json:"id"`
	Tweet *Tweet `json:"tweet"`
	Reaction string `json:"reaction"`
	User *User `json:"user"`

	CreatedAt int64 `json:"createdAt"`
}

func MakeTweet(userId int, text string, mediaPath string) (Tweet, error) {
	var tweet Tweet

	createdAt := time.Now().Unix()

	res, err := store.DB.Exec(
		"insert into tweets (userId, text, mediaPath, createdAt) values (?, ?, ?, ?)",
		userId, text, mediaPath, createdAt,
	)

	if err != nil {
		return tweet, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return tweet, err
	}

	// TODO: Do we want to fetch this at all?
	user, err := GetUserById(userId)

	if err != nil {
		return tweet, err
	}

	tweet.Id = int(id)
	tweet.CreatedAt = createdAt
	tweet.User = &user
	tweet.Text = text
	tweet.MediaPath = mediaPath

	return tweet, err
}

func GetTweetById(tweetId int) (Tweet, error) {
	var t Tweet

	t.Id = tweetId

	var mediaPath sql.NullString

	err := store.DB.
		QueryRow(`select text, mediaPath, createdAt from tweets where id = ?`, 
			tweetId,
		).Scan(&t.Text, &mediaPath, &t.CreatedAt)

	if (mediaPath.Valid) {
		t.MediaPath = mediaPath.String
	}

	// Not checking for `err` above because if it's
	// filled in, we can worry about that in the caller
	return t, err
}

func (t *Tweet) GetFullMediaPath() (string, error) {
	// TODO: Look for file
	path, _ := utils.GetStringOrDefault("MEDIA_PATH", "./assets/media")
	return fmt.Sprintf("%s/%s", path, t.MediaPath), nil
}

func (t *Tweet) MakeRetweet(userId int) (error) {
	_, err := store.DB.Exec(
		"insert into retweets (originalTweetId, userId) values (?, ?)",
		t.Id, userId,
	)

	return err
}

func (t *Tweet) DeleteRetweet(userId int) (error) {
	_, err := store.DB.Exec(
		"delete from retweets where userId = ? and tweetId = ?",
		userId, t.Id,
	)

	return err
}

func (t *Tweet) DeleteTweet() (error) {
	_, err := store.DB.Exec(
		"delete from tweets where id = ?",
		t.Id,
	)

	return err
}

func (t *Tweet) MakeReaction(userId int, reaction string) (error) {
	_, err := store.DB.Exec(
		"insert into reactions (originalTweetId, userId, reaction) values (?, ?, ?)",
		t.Id, userId, reaction,
	)

	return err
}

func (t *Tweet) DeleteReaction(userId int) (error) {
	_, err := store.DB.Exec(
		"delete from reactions where tweetId = ? and userId = ?",
		t.Id, userId,
	)

	return err
}

func (t *Tweet) ReactionCount() (map[string]int, error) {
	rows, err := store.DB.Query(
		"select reaction, count(id) from reactions where tweetId = ? group by reaction",
		t.Id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	reactionCount := make(map[string]int)

	for rows.Next() {
		var reaction string
		var count int
		err := rows.Scan(&reaction, &count)

		if err != nil {
			return nil, err
		}

		reactionCount[reaction] = count
	}

	return reactionCount, err
}

func (t *Tweet) RetweetCount() (int, error) {
	var count int

	err := store.DB.QueryRow(
		"select count(id) from retweets where tweetId = ?",
		t.Id,
	).Scan(&count)

	return count, err
}