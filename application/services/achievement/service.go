package service

import (
	"errors"
	achievement "github.com/scarlettmiss/engine-w/application/domain/achievement"
)

type Service struct {
	achievements achievement.Repository
}

func New(achievements achievement.Repository) (*Service, error) {
	if achievements == nil {
		return nil, errors.New("invalid achievements repo")
	}

	return &Service{
		achievements: achievements,
	}, nil
}

func (s *Service) CreateAchievement(name string) (*achievement.Achievement, error) {
	return s.achievements.CreateAchievement(name)
}

func (s *Service) Achievement(id string) (*achievement.Achievement, error) {
	return s.achievements.Achievement(id)
}
