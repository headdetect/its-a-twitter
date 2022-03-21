package models

import "time"

type Tweet struct {
	Id string
	User *User
	Text string
	MediaPath string

	CreatedAt time.Time
}

type Retweet struct {
	Id string
	Tweet *Tweet
	User *User

	createdAt time.Time
}

type Reaction struct {
	Id string
	Tweet *Tweet
	Reaction string
	User *User

	CreatedAt time.Time
}

type Timeline struct {
	User *User
	Tweets []*Tweet
}