package model

type User struct {
	Id int
	Username string
	DisplayName string
	ProfilePicPath string

	CreatedAt int64
}

type Follow struct {
	Id int

	LurkUser *User
	FollowedUser *User

	CreatedAt int64
}

