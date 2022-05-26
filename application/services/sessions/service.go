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

func (s *Service) CreateSession(userId string, capacity int, rating int, constraint session.Constraint) (*session.Session, error) {
	sess, err := s.sessions.CreateSession(capacity, rating, constraint)
	if err != nil {
		return nil, err
	}

	err = s.JoinSession(sess.Id(), userId)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *Service) JoinSession(id, userId string) error {
	sess, err := s.sessions.Session(id)
	if err != nil {
		return err
	}
	u, err := s.users.User(userId)
	if err != nil {
		return err
	}
	return sess.AddUser(u)
}

func (s *Service) LeaveSession(id, userId string) error {
	sess, err := s.sessions.Session(id)
	if err != nil {
		return err
	}

	err = sess.RemoveUser(userId)
	if err != nil {
		return err
	}

	if len(sess.Users()) == 0 {
		return s.sessions.DeleteSession(id)
	}

	return s.sessions.UpdateSession(sess)
}
