package service

import (
	"errors"
	"github.com/scarlettmiss/engine-w/application/domain/user"
)

type Service struct {
	users user.Repository
}

func New(users user.Repository) (*Service, error) {
	if users == nil {
		return nil, errors.New("invalid user repo")
	}

	return &Service{
		users: users,
	}, nil
}

func (s *Service) CreateUser(username string, password string) (*user.User, error) {
	return s.users.CreateUser(username, password)
}

func (s *Service) User(id string) (*user.User, error) {
	return s.users.User(id)
}

func (s *Service) Users() map[string]*user.User {
	return s.users.Users()
}

func (s *Service) UpdateUser(user *user.User) error {
	u, err := s.users.User(user.Id())
	if err != nil {
		return err
	}

	err = s.users.UpdateUser(u)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteUser(id string) error {
	return s.users.DeleteUser(id)
}

func (s *Service) Authenticate(username string, password string) (*user.User, error) {
	return s.users.CheckPassword(username, password)
}
