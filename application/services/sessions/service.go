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

func (s *Service) CreateSession(u *user.User, capacity int, minRating int, maxRating int, constraint string) (*session.Session, error) {
	constr, err := session.ParseConstraint(constraint)
	if err != nil {
		return nil, err
	}
	sess, err := s.sessions.CreateSession(u, capacity, minRating, maxRating, constr)
	if err != nil {
		return nil, err
	}

	_, err = s.JoinSession(sess.Id(), u.Id())
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *Service) UserSession(userId string) (*session.Session, error) {
	return s.sessions.UserSession(userId)
}

func (s *Service) JoinSession(id string, userId string) (*session.Session, error) {
	sess, err := s.sessions.Session(id)
	if err != nil {
		return nil, err
	}

	_, err = s.UserSession(userId)
	if err == nil {
		return nil, session.ErrUserInSession
	} else if !errors.Is(err, session.ErrNotFound) {
		return nil, err
	}

	if len(sess.Users()) == sess.Capacity() {
		return nil, session.ErrSessionFull
	}

	u, err := s.users.User(userId)
	if err != nil {
		return nil, err
	}

	err = sess.AddUser(u)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Service) RequestJoinSession(userId string) (*session.Session, error) {
	sessions := s.sessions.Sessions()

	if len(sessions) == 0 {
		return nil, session.ErrNoSessions
	}

	var sess *session.Session
	var err error
	for k := range sessions {
		sess, err = s.JoinSession(k, userId)
		if err == nil {
			break
		}
	}

	if err != nil && errors.Is(err, session.ErrSessionFull) {
		return nil, session.ErrNoSessions
	}

	return sess, err
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

func (s *Service) UpdateSession(sess *session.Session) error {
	return s.sessions.UpdateSession(sess)
}

func (s *Service) Sessions() map[string]*session.Session {
	return s.sessions.Sessions()
}

func (s *Service) Session(id string) (*session.Session, error) {
	return s.sessions.Session(id)
}
