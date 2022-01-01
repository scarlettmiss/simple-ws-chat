package user

type Repository interface {
	CreateUser() (*User, error)
	User() (*User, error)
	Users() (map[string]*User, error)
	UpdateUser(s *User) error
	DeleteUser(id string) error
}
