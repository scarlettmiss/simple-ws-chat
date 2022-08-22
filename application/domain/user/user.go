package user

import (
	"github.com/lithammer/shortuuid"
	"time"
)

type Achivement struct {
	Type      string
	Name      string
	CreatedOn time.Time
}

type User struct {
	id             string
	username       string
	password       string
	createdOn      time.Time
	updatedOn      time.Time
	lastSeenOnline time.Time
	online         bool
	friends        []string
	skillPoints    int
	achivement     map[string]*Achivement
	reputation     int
	deleted        bool
}

func New(username string, password string) *User {
	timestamp := time.Now()
	return &User{
		id:             shortuuid.New(),
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

func (u *User) Username() string {
	return u.username
}

func (u *User) SetUsername(username string) {
	u.username = username
}

func (u *User) CreatedOn() time.Time {
	return u.createdOn
}

func (u *User) SetCreatedOn(createdOn time.Time) {
	u.createdOn = createdOn
}

func (u *User) UpdatedOn() time.Time {
	return u.updatedOn
}

func (u *User) SetUpdatedOn(updatedOn time.Time) {
	u.updatedOn = updatedOn
}

func (u *User) LastSeenOnline() time.Time {
	return u.lastSeenOnline
}

func (u *User) SetLastSeenOnline(lastSeenOnline time.Time) {
	u.lastSeenOnline = lastSeenOnline
}

func (u *User) Online() bool {
	return u.online
}

func (u *User) SetOnline(online bool) {
	u.online = online
}

func (u *User) Friends() []string {
	return u.friends
}

func (u *User) SetFriends(friends []string) {
	u.friends = friends
}

func (u *User) SkillPoints() int {
	return u.skillPoints
}

func (u *User) SetSkillPoints(skillPoints int) {
	u.skillPoints = skillPoints
}

func (u *User) Achivement() map[string]*Achivement {
	return u.achivement
}

func (u *User) SetAchivement(achivement map[string]*Achivement) {
	u.achivement = achivement
}

func (u *User) Reputation() int {
	return u.reputation
}

func (u *User) SetReputation(reputation int) {
	u.reputation = reputation
}

func (u *User) Deleted() bool {
	return u.deleted
}

func (u *User) SetDeleted(deleted bool) {
	u.deleted = deleted
}

func (u *User) Password() string {
	return u.password
}

func (u *User) SetPassword(password string) {
	u.password = password
}
