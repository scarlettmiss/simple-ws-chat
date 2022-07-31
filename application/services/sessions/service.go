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

func (s *Service) CreateSession(userId string, capacity int, rating int, constraint string) (*session.Session, error) {
	constr, err := session.ParseConstraint(constraint)
	if err != nil {
		return nil, err
	}
	sess, err := s.sessions.CreateSession(userId, capacity, rating, constr)
	if err != nil {
		return nil, err
	}

	err = s.JoinSession(sess.Id(), userId)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *Service) UserSession(userId string) (*session.Session, error) {
	return s.sessions.UserSession(userId)
}

func (s *Service) JoinSession(id string, userId string) error {
	sess, err := s.sessions.Session(id)
	if err != nil {
		return err
	}

	_, err = s.UserSession(userId)
	if err == nil {
		return session.ErrUserInSession
	} else if !errors.Is(err, session.ErrNotFound) {
		return err
	}

	if len(sess.Users()) == sess.Capacity() {
		return errors.New("session is full")
	}

	u, err := s.users.User(userId)
	if err != nil {
		return err
	}
	return sess.AddUser(u)
}

func (s *Service) LeaveSession(id string, userId string) error {
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
