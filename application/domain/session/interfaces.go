package session

import (
	"errors"
	"github.com/scarlettmiss/engine-w/application/domain/user"
)

var (
	// ErrNotFound is returned when a session is not found
	ErrNotFound      = errors.New("session not found")
	ErrUserInSession = errors.New("user already in a session")
	ErrSessionFull   = errors.New("session is full")
	ErrNoSessions    = errors.New("no Session available")
)

type Service interface {
	CreateSession(owner *user.User, capacity int, minRating int, maxRating int, constraint string) (*Session, error)
	JoinSession(id string, userId string) (*Session, error)
	LeaveSession(id string, userId string) error
	UserSession(userId string) (*Session, error)
	Sessions() map[string]*Session
	Session(id string) (*Session, error)
	RequestJoinSession(userId string) (*Session, error)
	UpdateSession(s *Session) error
}

type Repository interface {
	CreateSession(owner *user.User, capacity int, minRating int, maxRating int, constraint Constraint) (*Session, error)
	Session(id string) (*Session, error)
	Sessions() map[string]*Session
	UpdateSession(s *Session) error
	DeleteSession(id string) error
	UserSession(userId string) (*Session, error)
}
