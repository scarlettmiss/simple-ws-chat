package session

import (
	"errors"
	"github.com/lithammer/shortuuid"
	"github.com/scarlettmiss/engine-w/application/domain/user"
)

type Session struct {
	id    string
	users map[string]*user.User
}

func New() *Session {
	return &Session{
		id:    shortuuid.New(),
		users: map[string]*user.User{},
	}
}

func (s *Session) Users() map[string]*user.User {
	return s.users
}

func (s *Session) Id() string {
	return s.id
}

func (s *Session) AddUser(u *user.User) error {
	_, exists := s.users[u.Id()]
	if exists {
		return errors.New("user already exists")
	}

	s.users[u.Id()] = u

	return nil
}

func (s *Session) RemoveUser(userId string) error {
	_, exists := s.users[userId]
	if !exists {
		return errors.New("user does not exist")
	}

	delete(s.users, userId)

	return nil
}
