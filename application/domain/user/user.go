package user

import (
	"github.com/lithammer/shortuuid"
	"github.com/scarlettmiss/engine-w/application/domain/achievement"
	"time"
)

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
	points         int
	gamesPlayed    int
	achievements   map[string]*achievement.Achievement
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
		achievements:   map[string]*achievement.Achievement{},
		skillPoints:    0, // xp
		points:         0, // xp
		gamesPlayed:    0,
		reputation:     100, // 0 to 100 worse to best depending on the player behavior
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

func (u *User) Achievements() map[string]*achievement.Achievement {
	return u.achievements
}

func (u *User) AddAchievement(achievement *achievement.Achievement) {
	u.achievements[achievement.Id()] = achievement
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

func (u *User) Points() int {
	return u.points
}

func (u *User) SetPoints(points int) {
	u.points = points
}

func (u *User) GamesPlayed() int {
	return u.gamesPlayed
}

func (u *User) SetGamesPlayed(gamesPlayed int) {
	u.gamesPlayed = gamesPlayed
}
