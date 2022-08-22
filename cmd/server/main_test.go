package main_test

import (
	"github.com/scarlettmiss/engine-w/application"
	"github.com/scarlettmiss/engine-w/application/domain/session"
	"github.com/scarlettmiss/engine-w/application/repositories/sessionrepo"
	"github.com/scarlettmiss/engine-w/application/repositories/userrepo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSessionJoin(t *testing.T) {
	sessionRepo := sessionrepo.New()
	userRepo := userrepo.New()

	app, err := application.New(sessionRepo, userRepo)
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

}
