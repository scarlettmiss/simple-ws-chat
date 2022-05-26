package session

import "errors"

var (
	// ErrNotFound is returned when a session is not found
	ErrNotFound = errors.New("session not found")
)

type Service interface {
	CreateSession(userId string, capacity int, rating int, constraint Constraint) (*Session, error)
	JoinSession(id, userId string) error
	LeaveSession(id, userId string) error
}

type Repository interface {
	CreateSession(capacity int, rating int, constraint Constraint) (*Session, error)
	Session(id string) (*Session, error)
	Sessions() (map[string]*Session, error)
	UpdateSession(s *Session) error
	DeleteSession(id string) error
}
