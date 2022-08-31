package achievement

import (
	"errors"
	"github.com/lithammer/shortuuid"
	"time"
)

var (
	// ErrNotFound is returned when a session is not found
	ErrNotFound = errors.New("achievement not found")
	ErrExists   = errors.New("achievement exists")
)

type Achievement struct {
	id        string
	name      string
	createdOn time.Time
}

func New(name string) *Achievement {
	timestamp := time.Now()
	return &Achievement{
		id:        shortuuid.New(),
		name:      name,
		createdOn: timestamp,
	}
}

func (a *Achievement) Id() string {
	return a.id
}

func (a *Achievement) Name() string {
	return a.name
}

func (a *Achievement) SetName(name string) {
	a.name = name
}

func (a *Achievement) CreatedOn() time.Time {
	return a.createdOn
}
