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

func (s *Service) CreateUser(username string, email string, password string) (*user.User, error) {
	return s.users.CreateUser(username, email, password)
}

func (s *Service) UpdateUser(userId string, username *string, email *string, password *string) (*user.User, error) {
	u, err := s.users.User(userId)
	if err != nil {
		return nil, err
	}

	if username != nil {
		u.Username = *username
	}
	if email != nil {
		u.Email = *email
	}
	if password != nil {
		u.Password = *password
	}

	err = s.users.UpdateUser(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) DeleteUser(id string) error {
	return s.users.DeleteUser(id)
}
