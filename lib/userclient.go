package lib

import "time"

type User struct {
	id             int32
	username       string
	password       string
	avatarFilename string
	createdAt      time.Time
	updatedAt      time.Time
}
