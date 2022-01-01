package service

import (
	"errors"
	"github.com/scarlettmiss/engine-w/application/domain/session"
	"github.com/scarlettmiss/engine-w/application/domain/user"
)

type Service struct {
	sessions session.Repository
	users    user.Repository
}

func (s *Service) CreateSession() (*session.Session, error) {
	panic("implement me")
}

func (s *Service) JoinSession(u user.User, sess *session.Session) error {
	panic("implement me")
}

func (s *Service) LeaveSession(id string) error {
	panic("implement me")
}

func New(sessions session.Repository, users user.Repository) (*Service, error) {
	if sessions == nil {
		return nil, errors.New("invalid session repo")
	}

	if users == nil {
		return nil, errors.New("invalid user repo")
	}

	return &Service{
		sessions: sessions,
		users:    users,
	}, nil
}
