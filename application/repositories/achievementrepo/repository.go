package achievementrepo

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
	_, err := r.AchievementByName(name)
	if err == nil {
		return nil, achievement.ErrExists
	} else if !errors.Is(err, achievement.ErrNotFound) {
		return nil, err
	}

	a := achievement.New(name)

	r.achievements[a.Id()] = a

	return a, nil
}

func (r *Repository) Achievement(id string) (*achievement.Achievement, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	a, ok := r.achievements[id]
	if !ok {
		return nil, achievement.ErrNotFound
	}

	return a, nil
}

func (r *Repository) AchievementByName(name string) (*achievement.Achievement, error) {
	for _, u := range r.achievements {
		if u.Name() == name {
			return u, nil
		}
	}
	return nil, achievement.ErrNotFound
}

func (r *Repository) Achievements() map[string]*achievement.Achievement {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.achievements
}
