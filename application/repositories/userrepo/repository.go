package userrepo

import (
	"errors"
	"github.com/scarlettmiss/engine-w/application/domain/user"
	"sync"
	"time"
)

type Repository struct {
	mux   sync.Mutex
	users map[string]*user.User
}

func New() *Repository {
	return &Repository{
		users: map[string]*user.User{},
	}
}

func (r *Repository) CreateUser(username string, password string) (*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, err := r.userByUsername(username)
	if err == nil {
		return nil, user.ErrUserExists
	} else if !errors.Is(err, user.ErrNotFound) {
		return nil, err
	}

	u := user.New(username, password)

	r.users[u.Id()] = u
	return u, nil
}

func (r *Repository) User(id string) (*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	u, ok := r.users[id]
	if !ok {
		return nil, user.ErrNotFound
	}

	return u, nil
}

func (r *Repository) userByUsername(username string) (*user.User, error) {
	for _, u := range r.users {
		if u.Username() == username {
			return u, nil
		}
	}
	return nil, user.ErrNotFound
}

func (r *Repository) UserByUsername(username string) (*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.userByUsername(username)
}

func (r *Repository) Users() (map[string]*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.users, nil
}

func (r *Repository) UpdateUser(u *user.User) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	u.SetUpdatedOn(time.Now())

	if !u.Online() {
		u.SetLastSeenOnline(time.Now())
	}

	_, ok := r.users[u.Id()]
	if !ok {
		return user.ErrNotFound
	}
	r.users[u.Id()] = u

	return nil

}

func (r *Repository) DeleteUser(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.users[id]
	if !ok {
		return user.ErrNotFound
	}

	delete(r.users, id)

	return nil

}

func (r *Repository) CheckPassword(username string, password string) (*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	u, err := r.userByUsername(username)
	if err != nil {
		return nil, user.ErrAuthentication
	}

	if u.Password() == password {
		return u, nil
	}

	return nil, user.ErrAuthentication
}
