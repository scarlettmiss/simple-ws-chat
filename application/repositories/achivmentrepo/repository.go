package sessionrepo

import (
	"errors"
	achievement "github.com/scarlettmiss/engine-w/application/domain/achievement"
	"sync"
)

type Repository struct {
	mux          sync.Mutex
	achievements map[string]*achievement.Achievement
}

func New() *Repository {
	return &Repository{
		achievements: map[string]*achievement.Achievement{},
	}
}

func (r *Repository) CreateAchievement(name string) (*achievement.Achievement, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	_, err := r.Achievement(name)
	if err == nil {
		return nil, achievement.ErrExists
	} else if !errors.Is(err, achievement.ErrNotFound) {
		return nil, err
	}

	a := achievement.New(name)

	r.achievements[a.Name()] = a

	return a, nil
}

func (r *Repository) Achievement(name string) (*achievement.Achievement, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	a, ok := r.achievements[name]
	if !ok {
		return nil, achievement.ErrNotFound
	}

	return a, nil
}

func (r *Repository) Achievements() map[string]*achievement.Achievement {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.achievements
}
