package application

import (
	"github.com/scarlettmiss/engine-w/application/domain/session"
	"github.com/scarlettmiss/engine-w/application/domain/user"
	sessionservice "github.com/scarlettmiss/engine-w/application/services/sessions"
)

type Application struct {
	sessionService session.Service
}

func New(sessions session.Repository, users user.Repository) (*Application, error) {
	ss, err := sessionservice.New(sessions, users)
	if err != nil {
		return nil, err
	}
	app := Application{sessionService: ss}

	return &app, nil
}

func (app *Application) CreateSession() (*session.Session, error) {
	sess, err := app.sessionService.CreateSession()
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (app *Application) JoinSession(user user.User, sess *session.Session) error {
	err := app.sessionService.JoinSession(user, sess)
	if err != nil {
		return err
	}
	return nil
}
