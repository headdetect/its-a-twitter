package model

import (
	"database/sql"
	"time"

	"github.com/headdetect/its-a-twitter/api/store"
)

type Tweet struct {
	Id        int    `json:"id"`
	User      User   `json:"user"`
	Text      string `json:"text"`
	MediaPath string `json:"mediaPath"` // TODO: Do we replace this with a dynamic hashed version of the id?

	CreatedAt int64 `json:"createdAt"`
}

type Retweet struct {
	Tweet Tweet `json:"tweet"`
	User  User  `json:"user"`

	CreatedAt int64 `json:"createdAt"`
}

type Reaction struct {
	Tweet    Tweet  `json:"tweet"`
	Reaction string `json:"reaction"`
	User     User   `json:"user"`

	CreatedAt int64 `json:"createdAt"`
}

// It's technically not a const, but we treat it as one //
var AllowedReactions = [...]string{"laugh", "party", "sad", "heart", "raisedHands", "shocked"}

func MakeTweet(user User, text string, mediaPath string) (Tweet, error) {
	var tweet Tweet

	createdAt := time.Now().Unix()

	var sqlMediaPath sql.NullString

	if mediaPath != "" {
		sqlMediaPath = sql.NullString{
			String: mediaPath,
			Valid:  true,
		}
	}

	res, err := store.DB.Exec(
		"insert into tweets (userId, text, mediaPath, createdAt) values (?, ?, ?, ?)",
		user.Id, text, sqlMediaPath, createdAt,
	)

	if err != nil {
		return tweet, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return tweet, err
	}

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

func GetTweetById(tweetId int) (t Tweet, err error) {
	t.Id = tweetId

	var mediaPath sql.NullString

	err = store.DB.
		QueryRow(`select text, mediaPath, createdAt from tweets where id = ?`,
			tweetId,
		).Scan(&t.Text, &mediaPath, &t.CreatedAt)

	if mediaPath.Valid {
		t.MediaPath = mediaPath.String
	}

	return
}

func (t *Tweet) MakeRetweet(userId int) (err error) {
	_, err = store.DB.Exec(
		"insert into retweets (tweetId, userId) values (?, ?)",
		t.Id, userId,
	)

	return
}

func (t *Tweet) DeleteRetweet(userId int) (err error) {
	_, err = store.DB.Exec(
		"delete from retweets where userId = ? and tweetId = ?",
		userId, t.Id,
	)

	return
}

func (t *Tweet) DeleteTweet() (err error) {
	_, err = store.DB.Exec(
		"delete from tweets where id = ?",
		t.Id,
	)

	return
}

func (t *Tweet) MakeReaction(userId int, reaction string) (err error) {
	_, err = store.DB.Exec(
		"insert into reactions (tweetId, userId, reaction) values (?, ?, ?)",
		t.Id, userId, reaction,
	)

	return
}

func (t *Tweet) DeleteReaction(userId int, reaction string) (err error) {
	_, err = store.DB.Exec(
		"delete from reactions where tweetId = ? and userId = ? and reaction = ?",
		t.Id, userId, reaction,
	)

	return
}

func (t *Tweet) ReactionCount() (map[string]int, error) {
	rows, err := store.DB.Query(
		"select reaction, count(*) from reactions where tweetId = ? group by reaction",
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

func (t *Tweet) RetweetCount() (count int, err error) {
	err = store.DB.QueryRow(
		"select count(*) from retweets where tweetId = ?",
		t.Id,
	).Scan(&count)

	return
}

func (t *Tweet) UserRetweeted(userId int) (bool, error) {
	var count int

	err := store.DB.QueryRow(
		"select count(*) from retweets where tweetId = ? and userId = ?",
		t.Id, userId,
	).Scan(&count)

	return count == 1, err
}

func (t *Tweet) UserReactions(userId int) ([]string, error) {
	reactions := make([]string, 0)

	rows, err := store.DB.Query(
		"select reaction from reactions where tweetId = ? and userId = ?",
		t.Id, userId,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var reaction string

		if err = rows.Scan(&reaction); err != nil {
			return nil, err
		}

		reactions = append(reactions, reaction)
	}

	return reactions, err
}
