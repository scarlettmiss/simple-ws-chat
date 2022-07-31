package session

import (
	"errors"
)

var (
	// ErrNotFound is returned when a session is not found
	ErrNotFound      = errors.New("session not found")
	ErrUserInSession = errors.New("user already in a session")
)

type Service interface {
	CreateSession(userId string, capacity int, rating int, constraint string) (*Session, error)
	JoinSession(id string, userId string) error
	LeaveSession(id string, userId string) error
	UserSession(userId string) (*Session, error)
}

type Repository interface {
	CreateSession(userId string, capacity int, rating int, constraint Constraint) (*Session, error)
	Session(id string) (*Session, error)
	Sessions() (map[string]*Session, error)
	UpdateSession(s *Session) error
	DeleteSession(id string) error
	UserSession(userId string) (*Session, error)
}
