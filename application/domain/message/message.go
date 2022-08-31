package message

import (
	"errors"
	"github.com/lithammer/shortuuid"
	"time"
)

var (
	// ErrNotFound is returned when a session is not found
	ErrNotFound = errors.New("messages not found")
)

type Message struct {
	id        string
	userId    string
	message   string
	createdOn time.Time
}

func New(userId string, message string) *Message {
	timestamp := time.Now()
	return &Message{
		id:        shortuuid.New(),
		userId:    userId,
		message:   message,
		createdOn: timestamp,
	}
}

func (m *Message) Id() string {
	return m.id
}

func (m *Message) UserId() string {
	return m.userId
}

func (m *Message) SetUserId(userId string) {
	m.userId = userId
}

func (m *Message) Message() string {
	return m.message
}

func (m *Message) SetMessage(message string) {
	m.message = message
}

func (m *Message) CreatedOn() time.Time {
	return m.createdOn
}
