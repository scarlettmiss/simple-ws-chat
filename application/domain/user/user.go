package user

import (
	"github.com/lithammer/shortuuid"
	"time"
)

type User struct {
	id             string
	Email          string
	Username       string
	Password       string
	createdOn      time.Time
	UpdatedOn      time.Time
	LastSeenOnline time.Time
	Online         bool
	Friends        []string
	SkillPoints    int
	Reputation     int
	Deleted        bool
}

func New(username string, email string, password string) *User {
	timestamp := time.Now()
	return &User{
		id:             shortuuid.New(),
		Email:          email,
		Username:       username,
		Password:       password,
		createdOn:      timestamp,
		UpdatedOn:      timestamp,
		LastSeenOnline: timestamp,
		Online:         true,
		Friends:        []string{},
		SkillPoints:    1000, // starting at 1000p so points won't be negative of the first game is a loss
		Reputation:     100,  // 0 to 100 worse to best depending on the player behavior
		Deleted:        false,
	}
}

func (u *User) Id() string {
	return u.id
}

func (u *User) update() string {
	return u.id
}
