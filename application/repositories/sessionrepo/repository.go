package sessionrepo

import (
	"errors"
	"github.com/scarlettmiss/engine-w/application/domain/session"
	"github.com/scarlettmiss/engine-w/application/domain/user"
	"sync"
)

type Repository struct {
	mux      sync.Mutex
	sessions map[string]*session.Session
}

func New() *Repository {
	return &Repository{
		sessions: map[string]*session.Session{},
	}
}

func (r *Repository) CreateSession(u *user.User, capacity int, minRating, maxRating int, constraint session.Constraint) (*session.Session, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	_, err := r.userSession(u.Id())
	if err == nil {
		return nil, session.ErrUserInSession
	} else if !errors.Is(err, session.ErrNotFound) {
		return nil, err
	}
	sess := session.New(u, capacity, minRating, maxRating, constraint)

	r.sessions[sess.Id()] = sess

	return sess, nil
}

func (r *Repository) Session(id string) (*session.Session, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	sess, ok := r.sessions[id]
	if !ok {
		return nil, session.ErrNotFound
	}

	return sess, nil
}

func (r *Repository) Sessions() map[string]*session.Session {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.sessions
}

func (r *Repository) userSession(userId string) (*session.Session, error) {
	for _, s := range r.sessions {
		if s.Users()[userId] != nil {
			return s, nil
		}
	}
	return nil, session.ErrNotFound
}

func (r *Repository) UserSession(userId string) (*session.Session, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	return r.userSession(userId)
}

func (r *Repository) UpdateSession(s *session.Session) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.sessions[s.Id()]
	if !ok {
		return session.ErrNotFound
	}

	r.sessions[s.Id()] = s

	return nil
}

func (r *Repository) DeleteSession(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, ok := r.sessions[id]
	if !ok {
		return session.ErrNotFound
	}

	delete(r.sessions, id)

	return nil
}
