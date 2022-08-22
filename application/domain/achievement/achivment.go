package achievment

import (
	"errors"
	"time"
)

var (
	// ErrNotFound is returned when a session is not found
	ErrNotFound = errors.New("session not found")
	ErrExists   = errors.New("achievement exists")
)

type Achievement struct {
	name      string
	createdOn time.Time
}

func New(name string) *Achievement {
	timestamp := time.Now()
	return &Achievement{
		name:      name,
		createdOn: timestamp,
	}
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

func (a *Achievement) SetCreatedOn(createdOn time.Time) {
	a.createdOn = createdOn
}
