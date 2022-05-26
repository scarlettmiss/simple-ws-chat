package userrepo

import (
	"github.com/scarlettmiss/engine-w/application/domain/user"
	"sync"
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

func (r *Repository) CreateUser(username string, email string, password string) (*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	u := user.New(username, email, password)

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

func (r *Repository) Users() (map[string]*user.User, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.users, nil
}

func (r *Repository) UpdateUser(u *user.User) error {
	r.mux.Lock()
	defer r.mux.Unlock()

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
