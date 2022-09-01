package movementrepo

import (
	"github.com/scarlettmiss/engine-w/application/domain/achievement"
	"github.com/scarlettmiss/engine-w/application/domain/movement"
	"sync"
)

type Repository struct {
	mux       sync.Mutex
	movements map[string]*movement.Movement
}

func New() *Repository {
	return &Repository{
		movements: map[string]*movement.Movement{},
	}
}

func (r *Repository) CreateMovement(movementType string, userId string, acceleration *movement.Point, position *movement.Point) (*movement.Movement, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	a := movement.New(movementType, userId, acceleration, position)

	r.movements[a.Id()] = a

	return a, nil
}

func (r *Repository) Movement(id string) (*movement.Movement, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	a, ok := r.movements[id]
	if !ok {
		return nil, achievement.ErrNotFound
	}

	return a, nil
}

func (r *Repository) MovementsByUserId(userId string) map[string]*movement.Movement {
	userMessages := map[string]*movement.Movement{}
	for _, v := range r.movements {
		if v.UserId() == userId {
			userMessages[v.Id()] = v
		}
	}
	return userMessages
}

func (r *Repository) Movements() map[string]*movement.Movement {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.movements
}
