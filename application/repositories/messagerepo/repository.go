package messagerepo

import (
	"github.com/scarlettmiss/engine-w/application/domain/achievement"
	"github.com/scarlettmiss/engine-w/application/domain/message"
	"sync"
)

type Repository struct {
	mux      sync.Mutex
	messages map[string]*message.Message
}

func New() *Repository {
	return &Repository{
		messages: map[string]*message.Message{},
	}
}

func (r *Repository) CreateMessage(userId string, chatMessage string) (*message.Message, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	a := message.New(userId, chatMessage)

	r.messages[a.Id()] = a

	return a, nil
}

func (r *Repository) Message(id string) (*message.Message, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	a, ok := r.messages[id]
	if !ok {
		return nil, achievement.ErrNotFound
	}

	return a, nil
}

func (r *Repository) MessagesByUserId(userId string) map[string]*message.Message {
	userMessages := map[string]*message.Message{}
	for _, u := range r.messages {
		if u.UserId() == userId {
			userMessages[u.Id()] = u
		}
	}
	return userMessages
}

func (r *Repository) Messages() map[string]*message.Message {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.messages
}
