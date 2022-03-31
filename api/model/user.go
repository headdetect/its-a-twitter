package model

import (
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

	LurkUser *User
	FollowedUser *User

	CreatedAt int64
}

func GetUserWithPassByUsername(username string) (*User, string, error) {
	var user User
	var hashedPassword string

	err := store.DB.QueryRow(
			"select id, username, password, createdAt from users where username = ? limit 1",
			username,
		).Scan(
			&user.Id, &user.Username, &hashedPassword, &user.CreatedAt,
		)

	return &user, hashedPassword, err
}

func GetUserById(id int) (*User, error) {
	var user User

	err := store.DB.QueryRow(
			"select id, username, createdAt from user where id = ? limit 1", 
			id,
		).Scan(
			&user.Id, &user.Username, &user.CreatedAt,
		)

	return &user, err
}

func MakeUser(username string, passwordHash string) (*User, error) {
	var user User
	
	createdAt := time.Now().Unix()

	res, err := store.DB.Exec(
		"insert into user (username, password, createdAt) values (?, ?, ?, ?)",
		username, passwordHash, createdAt,
	)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return nil, err
	}

	user.Id = int(id)
	user.CreatedAt = createdAt
	user.Username = username

	return &user, err
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


func (u *User) FollowUser(userId int) (error) {
	_, err := store.DB.Exec(
		"insert into follows (lurkerUserId, followedUserId) values (?, ?)",
		u.Id, userId,
	)

	return err
}

func (u *User) UnFollowUser(userId int) (error) {
	_, err := store.DB.Exec(
		"delete from follows where lurkerUserId = ? and followedUserId = ?",
		u.Id, userId,
	)

	return err
}

func (u *User) GetFollowers() ([]*User, error) {
	rows, err := store.DB.
		Query(
			`select 
				followers.id,
				followers.username,
				followers.createdAt
			from follows f
			join users followers on f.lurkerUserId = followers.id
			where f.followedUserId = ?`, 
			u.Id,
		)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var followers []*User
	
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Username, &u.CreatedAt)

		if err != nil {
			return nil, err
		}

		followers = append(followers, &u)
	}

	return followers, nil
}

func (u *User) GetFollowing() ([]*User, error) {
	rows, err := store.DB.
		Query(
			`select 
				following.id,
				following.username,
				following.createdAt
			from follows
			join users following on follows.followedUserId = following.id
			where follows.lurkerUserId = ?`, 
			u.Id,
		)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var followering []*User
	
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Username, &u.CreatedAt)

		if err != nil {
			return nil, err
		}

		followering = append(followering, &u)
	}

	return followering, nil
}


func (u *User) GetTweets() ([]Tweet, error) {
	rows, err := store.DB.
		Query(`select text, mediaPath, createdAt from tweets where userId = ?`, 
			u.Id,
		)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tweets []Tweet
	
	for rows.Next() {
		var t Tweet
		err := rows.Scan(&t.Text, &t.MediaPath, &t.CreatedAt)

		if err != nil {
			return nil, err
		}

		t.User = u

		tweets = append(tweets, t)
	}

	return tweets, nil
}