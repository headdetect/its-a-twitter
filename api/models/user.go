package models

import "time"

type User struct {
	Username string
	DisplayName string
	HashedPassword string
	ProfilePicPath string

	CreatedAt time.Time
}

