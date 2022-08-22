package application

import (
	"github.com/scarlettmiss/engine-w/application/domain/session"
	"github.com/scarlettmiss/engine-w/application/domain/user"
	sessionservice "github.com/scarlettmiss/engine-w/application/services/sessions"
	userservice "github.com/scarlettmiss/engine-w/application/services/user"
)

type Application struct {
	sessionService session.Service
	userService    user.Service
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
		userService:    us,
	}

	return &app, nil
}

func (app *Application) CreateSession(owner *user.User, capacity int, minRating int, maxRating int, constraint string) (*session.Session, error) {
	return app.sessionService.CreateSession(owner, capacity, minRating, maxRating, constraint)
}

func (app *Application) JoinSession(id string, userId string) (*session.Session, error) {
	return app.sessionService.JoinSession(id, userId)
}

func (app *Application) LeaveSession(id, userId string) error {
	return app.sessionService.LeaveSession(id, userId)
}

func (app *Application) UserSession(userId string) (*session.Session, error) {
	return app.sessionService.UserSession(userId)
}

func (app *Application) GetSessions() map[string]*session.Session {
	return app.sessionService.Sessions()
}

func (app *Application) CreateUser(username string, password string) (*user.User, error) {
	return app.userService.CreateUser(username, password)
}

func (app *Application) Authenticate(username string, password string) (*user.User, error) {
	return app.userService.Authenticate(username, password)
}

func (app *Application) UpdateUser(user *user.User) error {
	return app.userService.UpdateUser(user)
}

func (app *Application) User(userId string) (*user.User, error) {
	return app.userService.User(userId)
}

func (app *Application) DeleteUser(id string) error {
	return app.userService.DeleteUser(id)
}

func (app *Application) RequestJoinSession(userId string) (*session.Session, error) {
	return app.sessionService.RequestJoinSession(userId)
}
