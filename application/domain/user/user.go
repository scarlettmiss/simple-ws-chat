package user

import (
	"github.com/lithammer/shortuuid"
	"time"
)

type User struct {
	Id             string
	Username       string
	password       string
	CreatedOn      time.Time
	UpdatedOn      time.Time
	LastSeenOnline time.Time
	Online         bool
	Friends        []string
	SkillPoints    int
	Reputation     int
	Deleted        bool
}

func New(username string, password string) *User {
	timestamp := time.Now()
	return &User{
		Id:             shortuuid.New(),
		Username:       username,
		password:       password,
		CreatedOn:      timestamp,
		UpdatedOn:      timestamp,
		LastSeenOnline: timestamp,
		Online:         true,
		Friends:        []string{},
		SkillPoints:    1000, // starting at 1000p so points won't be negative of the first game is a loss
		Reputation:     100,  // 0 to 100 worse to best depending on the player behavior
		Deleted:        false,
	}
}

func (u *User) Password() string {
	return u.password
}

func (u *User) SetPassword(password string) {
	u.password = password
}

func (u *User) update() string {
	return u.Id
}
