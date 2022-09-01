package application

import (
	achievement "github.com/scarlettmiss/engine-w/application/domain/achievement"
	"github.com/scarlettmiss/engine-w/application/domain/message"
	"github.com/scarlettmiss/engine-w/application/domain/movement"
	"github.com/scarlettmiss/engine-w/application/domain/session"
	"github.com/scarlettmiss/engine-w/application/domain/user"
	achievementService "github.com/scarlettmiss/engine-w/application/services/achievement"
	messagesService "github.com/scarlettmiss/engine-w/application/services/message"
	movementService "github.com/scarlettmiss/engine-w/application/services/movement"
	sessionservice "github.com/scarlettmiss/engine-w/application/services/sessions"
	userservice "github.com/scarlettmiss/engine-w/application/services/user"
)

type Application struct {
	sessionService     session.Service
	userService        user.Service
	achievementService achievement.Service
	messageService     message.Service
	movementsService   movement.Service
}

func New(sessions session.Repository, users user.Repository, achievements achievement.Repository, messages message.Repository, movement movement.Repository) (*Application, error) {
	ss, err := sessionservice.New(sessions, users)
	if err != nil {
		return nil, err
	}

	us, err := userservice.New(users)
	if err != nil {
		return nil, err
	}

	as, err := achievementService.New(achievements)
	if err != nil {
		return nil, err
	}

	ms, err := messagesService.New(messages)
	if err != nil {
		return nil, err
	}

	moveS, err := movementService.New(movement)
	if err != nil {
		return nil, err
	}

	app := Application{
		sessionService:     ss,
		userService:        us,
		achievementService: as,
		messageService:     ms,
		movementsService:   moveS,
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

func (app *Application) User(userId string) (*user.User, error) {
	return app.userService.User(userId)
}

func (app *Application) Users() map[string]*user.User {
	return app.userService.Users()
}

func (app *Application) DeleteUser(id string) error {
	return app.userService.DeleteUser(id)
}

func (app *Application) RequestJoinSession(userId string) (*session.Session, error) {
	return app.sessionService.RequestJoinSession(userId)
}

func (app *Application) Achievement(id string) (*achievement.Achievement, error) {
	return app.achievementService.Achievement(id)
}

func (app *Application) CreateAchievement(name string) (*achievement.Achievement, error) {
	return app.achievementService.CreateAchievement(name)
}

func (app *Application) CreateMessage(userId string, message string) (*message.Message, error) {
	return app.messageService.CreateMessage(userId, message)
}

func (app *Application) Messages() map[string]*message.Message {
	return app.messageService.Messages()
}

func (app *Application) SessionAddMessage(sessId string, userId string, message string) (*message.Message, error) {
	sess, err := app.sessionService.Session(sessId)
	if err != nil {
		return nil, err
	}
	m, err := app.CreateMessage(userId, message)
	if err != nil {
		return nil, err
	}
	err = sess.AddMessage(m)
	if err != nil {
		return nil, err
	}

	err = app.sessionService.UpdateSession(sess)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (app *Application) UserSetOnline(userId string, value bool) error {
	u, err := app.User(userId)
	if err != nil {
		return err
	}
	u.SetOnline(value)
	err = app.userService.UpdateUser(u)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) UserSetAccountInfo(userId string, username *string, password *string) error {
	u, err := app.User(userId)
	if err != nil {
		return err
	}

	if username != nil {
		u.SetUsername(*username)
	}

	if password != nil {
		u.SetPassword(*password)
	}

	err = app.userService.UpdateUser(u)
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) UserAddAchievement(userId string, achievementId string) error {
	a, err := app.Achievement(achievementId)
	if err != nil {
		return err
	}
	u, err := app.User(userId)
	if err != nil {
		return err
	}
	u.AddAchievement(a)
	err = app.userService.UpdateUser(u)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) CreateMovement(userId string, movementType string, acceleration *movement.Point, position *movement.Point) (*movement.Movement, error) {
	return app.movementsService.CreateMovement(movementType, userId, acceleration, position)
}

func (app *Application) Movements() map[string]*movement.Movement {
	return app.movementsService.Movements()
}

func (app *Application) SessionAddMovement(sessId string, userId string, movementType string, acceleration *movement.Point, position *movement.Point) (*movement.Movement, error) {
	sess, err := app.sessionService.Session(sessId)
	if err != nil {
		return nil, err
	}
	m, err := app.CreateMovement(userId, movementType, acceleration, position)
	if err != nil {
		return nil, err
	}
	err = sess.AddMovement(m)
	if err != nil {
		return nil, err
	}

	err = app.sessionService.UpdateSession(sess)
	if err != nil {
		return nil, err
	}
	return m, nil
}
