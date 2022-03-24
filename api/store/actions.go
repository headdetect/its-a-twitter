package store

import "github.com/headdetect/its-a-twitter/api/model"

func InsertFollowLink(lurkerUserId int, followedUserId int) (error) {
	return nil
}

func GetFollowers(userId int) ([]*model.User, error) {
	var user model.User

	err := db.
		QueryRow(
			"select id, username, displayName, profilePicHash, createdAt from user where id = ? limit 1", 
			userId,
		).
		Scan(&user.Id, &user.Username, &user.DisplayName, &user.ProfilePicPath, &user.CreatedAt)

	return nil, err
}

func InsertTweet(userId, tweet string, ) (error) {
	return nil
}

func GetTweetsByUser(userId int)

func GetUserWithPassByUsername(username string) (*model.User, string, error) {
	var user model.User
	var hashedPassword string

	err := db.
		QueryRow(
			"select id, username, displayName, profilePicHash, password, createdAt from users where username = ? limit 1",
			username,
		).
		Scan(&user.Id, &user.Username, &user.DisplayName, &user.ProfilePicPath, &hashedPassword, &user.CreatedAt)

	return &user, hashedPassword, err
}

func GetUserById(id int) (*model.User, error) {
	var user model.User

	err := db.
		QueryRow(
			"select id, username, displayName, profilePicHash, createdAt from user where id = ? limit 1", 
			id,
		).
		Scan(&user.Id, &user.Username, &user.DisplayName, &user.ProfilePicPath, &user.CreatedAt)

	return &user, err
}