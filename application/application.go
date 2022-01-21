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

func (app *Application) CreateSession(userId string) (*session.Session, error) {
	sess, err := app.sessionService.CreateSession(userId)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (app *Application) JoinSession(id, userId string) error {
	err := app.sessionService.JoinSession(id, userId)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) LeaveSession(id, userId string) error {
	err := app.sessionService.LeaveSession(id, userId)
	if err != nil {
		return err
	}
	return nil
}
