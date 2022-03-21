package models

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

	createdAt int64
}

type Reaction struct {
	Id int
	Tweet *Tweet
	Reaction string
	User *User

	CreatedAt int64
}

type Timeline struct {
	User *User
	Tweets []*Tweet
}