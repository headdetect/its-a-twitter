package model

import (
	"fmt"
	"time"

	"github.com/headdetect/its-a-twitter/api/store"
)

var MEDIA_PATH = "./assets/media"


type Tweet struct {
	Id int
	User *User
	Text string
	MediaPath string

	CreatedAt int64
}

type Retweet struct {
	Id int
	Tweet *Tweet
	User *User

	CreatedAt int64
}

type Reaction struct {
	Id int
	Tweet *Tweet
	Reaction string
	User *User

	CreatedAt int64
}


type GeneratedTimeline struct {
	User *User
	Tweets []*Tweet
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
	tweet.User = user
	tweet.Text = text
	tweet.MediaPath = mediaPath

	return tweet, err
}

func GetTweetById(tweetId int) (Tweet, error) {
	return Tweet{}, nil
}

func (t *Tweet) GetMediaPath() (string, error) {
	// TODO: Look for file

	return fmt.Sprintf("%s/%s", MEDIA_PATH, t.MediaPath), nil
}

func (t *Tweet) MakeRetweet(userId int) (error) {
	_, err := store.DB.Exec(
		"insert into retweets (originalTweetId, userId) values (?, ?)",
		t.Id, userId,
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


func (t *Tweet) RemoveReaction(userId int) (error) {
	_, err := store.DB.Exec(
		"delete from reactions where tweetId = ? and userId = ?",
		t.Id, userId,
	)

	return err
}
