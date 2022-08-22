package user

import "errors"

var (
	// ErrNotFound is returned when a user is not found
	ErrNotFound       = errors.New("user not found")
	ErrUserExists     = errors.New("user already exists")
	ErrAuthentication = errors.New("wrong credentials")
)

type Service interface {
	User(id string) (*User, error)
	CreateUser(username string, password string) (*User, error)
	Authenticate(username string, password string) (*User, error)
	UpdateUser(u *User) error
	DeleteUser(id string) error
}

type Repository interface {
	CreateUser(username string, password string) (*User, error)
	User(id string) (*User, error)
	UserByUsername(username string) (*User, error)
	Users() (map[string]*User, error)
	UpdateUser(u *User) error
	DeleteUser(id string) error
	CheckPassword(username string, password string) (*User, error)
}
