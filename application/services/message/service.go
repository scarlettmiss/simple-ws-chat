package service

import (
	"errors"
	"github.com/scarlettmiss/engine-w/application/domain/message"
)

type Service struct {
	messages message.Repository
}

func New(messages message.Repository) (*Service, error) {
	if messages == nil {
		return nil, errors.New("invalid messages repo")
	}

	return &Service{
		messages: messages,
	}, nil
}

func (s *Service) CreateMessage(userId string, chatMessage string) (*message.Message, error) {
	return s.messages.CreateMessage(userId, chatMessage)
}

func (s *Service) Messages() map[string]*message.Message {
	return s.messages.Messages()
}

func (s *Service) MessagesByUserId(userId string) map[string]*message.Message {
	return s.messages.MessagesByUserId(userId)
}
