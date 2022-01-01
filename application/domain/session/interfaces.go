package session

import "github.com/scarlettmiss/engine-w/application/domain/user"

type Service interface {
	CreateSession() (*Session, error)
	JoinSession(u user.User, s *Session) error
	LeaveSession(id string) error
}

type Repository interface {
	CreateSession() (*Session, error)
	Session(id string) (*Session, error)
	Sessions() (map[string]*Session, error)
	UpdateSession(s *Session) (*Session, error)
	DeleteSession(id string) error
}
