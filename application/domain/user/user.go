package user

import "github.com/lithammer/shortuuid"

type User struct {
	id string
}

func New() *User {
	return &User{
		id: shortuuid.New(),
	}
}

func (u *User) Id() string {
	return u.id
}
