package user

import (
	"github.com/lithammer/shortuuid"
	"time"
)

type User struct {
	id             string
	email          string
	username       string
	password       string
	createdOn      time.Time
	updatedOn      time.Time
	lastSeenOnline time.Time
	online         bool
	friends        []string
	skillPoints    int
	reputation     int
	deleted        bool
}

func New(username string, email string, password string) *User {
	timestamp := time.Now()
	return &User{
		id:             shortuuid.New(),
		email:          email,
		username:       username,
		password:       password,
		createdOn:      timestamp,
		updatedOn:      timestamp,
		lastSeenOnline: timestamp,
		online:         true,
		friends:        []string{},
		skillPoints:    1000, // starting at 1000p so points won't be negative of the first game is a loss
		reputation:     100,  // 0 to 100 worse to best depending on the player behavior
		deleted:        false,
	}
}

func (u *User) Id() string {
	return u.id
}

func (u *User) update() string {
	return u.id
}
