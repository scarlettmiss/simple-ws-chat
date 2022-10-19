package userrepo

import (
	"github.com/jmoiron/sqlx"
	dbUser "github.com/scarlettmiss/engine-w/application/db/user"
	"sync"
)

type Repository struct {
	mux      sync.Mutex
	db       *sqlx.DB
	dbSchema string
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db:       db,
		dbSchema: dbUser.Schema,
	}
}

func (r *Repository) Init() {
	r.db.MustExec(r.dbSchema)
}

//
//func (r *Repository) CreateUser(username string, password string) (*user.User, error) {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	_, err := r.userByUsername(username)
//	if err == nil {
//		return nil, user.ErrUserExists
//	} else if !errors.Is(err, user.ErrNotFound) {
//		return nil, err
//	}
//
//	u := user.New(username, password)
//
//	tx := r.db.MustBegin()
//	tx.MustExec("INSERT INTO user (id, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
//	tx.Commit()
//}
//
//func (r *Repository) User(id string) (*user.User, error) {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	u, ok := r.db[id]
//	if !ok {
//		return nil, user.ErrNotFound
//	}
//
//	return u, nil
//}
//
//func (r *Repository) userByUsername(username string) (*user.User, error) {
//	for _, u := range r.db {
//		if u.Username() == username {
//			return u, nil
//		}
//	}
//	return nil, user.ErrNotFound
//}
//
//func (r *Repository) UserByUsername(username string) (*user.User, error) {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	return r.userByUsername(username)
//}
//
//func (r *Repository) Users() map[string]*user.User {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	return r.db
//}
//
//func (r *Repository) UpdateUser(u *user.User) error {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	u.SetUpdatedOn(time.Now())
//
//	if !u.Online() {
//		u.SetLastSeenOnline(time.Now())
//	}
//
//	_, ok := r.db[u.Id()]
//	if !ok {
//		return user.ErrNotFound
//	}
//	r.db[u.Id()] = u
//
//	return nil
//
//}
//
//func (r *Repository) DeleteUser(id string) error {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	_, ok := r.db[id]
//	if !ok {
//		return user.ErrNotFound
//	}
//
//	delete(r.db, id)
//
//	return nil
//
//}
//
//func (r *Repository) CheckPassword(username string, password string) (*user.User, error) {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	u, err := r.userByUsername(username)
//	if err != nil {
//		return nil, user.ErrAuthentication
//	}
//
//	if u.Password() == password {
//		return u, nil
//	}
//
//	return nil, user.ErrAuthentication
//}
