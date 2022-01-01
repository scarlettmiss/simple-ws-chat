package user

type User struct {
	id string
}

func (u *User) Id() string {
	return u.id
}
