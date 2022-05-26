package user

import "errors"

var (
	// ErrNotFound is returned when a user is not found
	ErrNotFound = errors.New("user not found")
)

type Service interface {
	CreateUser(username string, email string, password string) (*User, error)
	UpdateUser(u *User) error
	DeleteUser(id string) error
}

type Repository interface {
	CreateUser(username string, email string, password string) (*User, error)
	User(id string) (*User, error)
	Users() (map[string]*User, error)
	UpdateUser(u *User) error
	DeleteUser(id string) error
}
