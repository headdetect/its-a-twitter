package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/headdetect/its-a-twitter/api/store"
)

var PROFILE_PIC_PATH = "./assets/profile"

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ProfilePicPath string `json:"profilePicPath"`

	CreatedAt int64 `json:"createdAt"`
}

type Follow struct {
	User         *User `json:"user"`
	FollowedUser *User `json:"followedUser"`

	CreatedAt int64 `json:"createdAt"`
}

func GetUserByUsernameWithPass(username string) (user User, hashedPassword string, email string, err error) {
	var profilePicPath sql.NullString

	err = store.DB.QueryRow(
		"select id, username, profilePicPath, email, password, createdAt from users where username = ? limit 1",
		username,
	).Scan(
		&user.Id, &user.Username, &profilePicPath, &email, &hashedPassword, &user.CreatedAt,
	)

	if profilePicPath.Valid {
		user.ProfilePicPath = profilePicPath.String
	}

	return
}

func GetUserById(id int) (user User, email string, err error) {
	var profilePicPath sql.NullString

	err = store.DB.QueryRow(
		"select id, username, profilePicPath, email, createdAt from users where id = ? limit 1",
		id,
	).Scan(
		&user.Id, &user.Username, &profilePicPath, &email, &user.CreatedAt,
	)

	if profilePicPath.Valid {
		user.ProfilePicPath = profilePicPath.String
	}

	return
}

func GetUserByUsername(username string) (user User, email string, err error) {
	var profilePicPath sql.NullString

	err = store.DB.QueryRow(
		"select id, username, profilePicPath, createdAt from users where username = ? limit 1",
		username,
	).Scan(
		&user.Id, &user.Username, &profilePicPath, &user.CreatedAt,
	)

	if profilePicPath.Valid {
		user.ProfilePicPath = profilePicPath.String
	}

	return
}

func GetUserByTweetId(tweetId int) (user User, email string, err error) {
	var profilePicPath sql.NullString

	err = store.DB.QueryRow(
		"select u.id, u.username, u.profilePicPath, u.email, u.createdAt from users u, tweets t where u.id = t.userId and t.id = ? limit 1",
		tweetId,
	).Scan(
		&user.Id, &user.Username, &profilePicPath, &user.Email, &user.CreatedAt,
	)

	if profilePicPath.Valid {
		user.ProfilePicPath = profilePicPath.String
	}

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


func (u *User) UpdateProfilePicPath(profilePicPath string) (err error) {
	_, err = store.DB.Exec(
		"update users set profilePicPath = ? where id = ? ",
		profilePicPath, u.Id,
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
