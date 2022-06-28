package application

import (
	"github.com/scarlettmiss/engine-w/application/domain/session"
	"github.com/scarlettmiss/engine-w/application/domain/user"
	sessionservice "github.com/scarlettmiss/engine-w/application/services/sessions"
	userservice "github.com/scarlettmiss/engine-w/application/services/user"
)

type Application struct {
	sessionService session.Service
	UserService    user.Service
}

func New(sessions session.Repository, users user.Repository) (*Application, error) {
	ss, err := sessionservice.New(sessions, users)
	if err != nil {
		return nil, err
	}

	us, err := userservice.New(users)
	if err != nil {
		return nil, err
	}

	app := Application{
		sessionService: ss,
		UserService:    us,
	}

	return &app, nil
}

func (app *Application) CreateSession(userId string, capacity int, rating int, constraint session.Constraint) (*session.Session, error) {
	return app.sessionService.CreateSession(userId, capacity, rating, constraint)
}

func (app *Application) JoinSession(id, userId string) error {
	return app.sessionService.JoinSession(id, userId)
}

func (app *Application) LeaveSession(id, userId string) error {
	return app.sessionService.LeaveSession(id, userId)
}

func (app *Application) CreateUser(username string, email string, password string) (*user.User, error) {
	return app.UserService.CreateUser(username, email, password)
}

func (app *Application) UpdateUser(userId string, username *string, email *string, password *string) (*user.User, error) {
	return app.UserService.UpdateUser(userId, username, email, password)
}

func (app *Application) DeleteUser(id string) error {
	return app.UserService.DeleteUser(id)
}
