package main_test

import (
	"github.com/scarlettmiss/engine-w/application"
	"github.com/scarlettmiss/engine-w/application/domain/session"
	"github.com/scarlettmiss/engine-w/application/repositories/achievementrepo"
	"github.com/scarlettmiss/engine-w/application/repositories/messagerepo"
	"github.com/scarlettmiss/engine-w/application/repositories/sessionrepo"
	"github.com/scarlettmiss/engine-w/application/repositories/userrepo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserCreation(t *testing.T) {
	sessionRepo := sessionrepo.New()
	userRepo := userrepo.New()
	achievementRepo := achievementrepo.New()
	messageRepo := messagerepo.New()

	app, err := application.New(sessionRepo, userRepo, achievementRepo, messageRepo)
	assert.Nil(t, err)

	u, err := app.CreateUser("scarlettmiss", "1234")
	assert.Nil(t, err)
	assert.Equal(t, u.Username(), "scarlettmiss")
}

func TestUserAuthentication(t *testing.T) {
	sessionRepo := sessionrepo.New()
	userRepo := userrepo.New()
	achievementRepo := achievementrepo.New()
	messageRepo := messagerepo.New()

	app, err := application.New(sessionRepo, userRepo, achievementRepo, messageRepo)
	assert.Nil(t, err)

	_, err = app.CreateUser("scarlettmiss", "1234")
	assert.Nil(t, err)

	u, err := app.Authenticate("scarlettmiss", "1234")
	assert.Nil(t, err)

	assert.Equal(t, u.Username(), "scarlettmiss")
}

func TestAchivementCreation(t *testing.T) {
	sessionRepo := sessionrepo.New()
	userRepo := userrepo.New()
	achievementRepo := achievementrepo.New()
	messageRepo := messagerepo.New()

	app, err := application.New(sessionRepo, userRepo, achievementRepo, messageRepo)
	assert.Nil(t, err)

	_, err = app.CreateAchievement("1st_Log")
	assert.Nil(t, err)
}

func TestMessageCreation(t *testing.T) {
	sessionRepo := sessionrepo.New()
	userRepo := userrepo.New()
	achievementRepo := achievementrepo.New()
	messageRepo := messagerepo.New()

	app, err := application.New(sessionRepo, userRepo, achievementRepo, messageRepo)
	assert.Nil(t, err)

	u, err := app.CreateUser("scarlettmiss", "1234")
	assert.Nil(t, err)

	_, err = app.CreateMessage(u.Id(), "1st_Log")
	assert.Nil(t, err)
}

func TestUserUpdate(t *testing.T) {
	sessionRepo := sessionrepo.New()
	userRepo := userrepo.New()
	achievementRepo := achievementrepo.New()
	messageRepo := messagerepo.New()

	app, err := application.New(sessionRepo, userRepo, achievementRepo, messageRepo)
	assert.Nil(t, err)

	u, err := app.CreateUser("scarlettmiss", "1234")
	assert.Nil(t, err)

	a, err := app.CreateAchievement("1st_Log")
	assert.Nil(t, err)

	err = app.UserAddAchievement(u.Id(), a.Id())
	assert.Nil(t, err)

	name := "sm"
	err = app.UserSetAccountInfo(u.Id(), &name, nil)
	assert.Nil(t, err)

	err = app.UserSetOnline(u.Id(), true)
	assert.Nil(t, err)

	assert.Equal(t, u.Online(), true)
	assert.Equal(t, u.Username(), "sm")
	assert.Len(t, u.Achievements(), 1)
	assert.Equal(t, u.Achievements()[a.Id()], a)

}

func TestSessionUpdate(t *testing.T) {
	sessionRepo := sessionrepo.New()
	userRepo := userrepo.New()
	achievementRepo := achievementrepo.New()
	messageRepo := messagerepo.New()

	app, err := application.New(sessionRepo, userRepo, achievementRepo, messageRepo)
	assert.Nil(t, err)

	u, err := app.CreateUser("scarlettmiss", "1234")
	assert.Nil(t, err)

	s, err := app.CreateSession(u, 4, 1000, 2000, "none")
	assert.Nil(t, err)
	assert.Len(t, s.Users(), 1)

	m, err := app.SessionAddMessage(s.Id(), u.Id(), "Hello")
	assert.Nil(t, err)
	assert.Len(t, s.Messages(), 1)
	assert.Equal(t, s.Messages()[0], m)
}

func TestSessionJoin(t *testing.T) {
	sessionRepo := sessionrepo.New()
	userRepo := userrepo.New()
	achievementRepo := achievementrepo.New()
	messageRepo := messagerepo.New()

	app, err := application.New(sessionRepo, userRepo, achievementRepo, messageRepo)
	assert.Nil(t, err)

	u, err := app.CreateUser("scarlettmiss", "1234")
	assert.Nil(t, err)

	s, err := app.CreateSession(u, 4, 1000, 2000, "none")
	assert.Nil(t, err)
	assert.Len(t, s.Users(), 1)

	err = app.LeaveSession(s.Id(), u.Id())
	assert.Nil(t, err)

	_, err = sessionRepo.Session(s.Id())
	assert.Equal(t, err, session.ErrNotFound)

	ss := app.GetSessions()
	assert.Len(t, ss, 0)
}
