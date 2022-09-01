package service

import (
	"errors"
	"github.com/scarlettmiss/engine-w/application/domain/movement"
)

type Service struct {
	movements movement.Repository
}

func New(movements movement.Repository) (*Service, error) {
	if movements == nil {
		return nil, errors.New("invalid messages repo")
	}

	return &Service{
		movements: movements,
	}, nil
}

func (s *Service) CreateMovement(movementType string, userId string, acceleration *movement.Point, position *movement.Point) (*movement.Movement, error) {
	return s.movements.CreateMovement(movementType, userId, acceleration, position)
}

func (s *Service) Movements() map[string]*movement.Movement {
	return s.movements.Movements()
}

func (s *Service) MovementsByUserId(userId string) map[string]*movement.Movement {
	return s.movements.MovementsByUserId(userId)
}
