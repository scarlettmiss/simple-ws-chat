package session

import (
	"errors"
	"github.com/lithammer/shortuuid"
	"github.com/scarlettmiss/engine-w/application/domain/message"
	"github.com/scarlettmiss/engine-w/application/domain/movement"
	"github.com/scarlettmiss/engine-w/application/domain/user"
)

type Constraint string

const (
	ConstraintNone           Constraint = "none"
	ConstraintFriendsOnly    Constraint = "friendsOnly"
	ConstraintInvitationOnly Constraint = "invitationOnly"
)

func ParseConstraint(s string) (Constraint, error) {
	switch Constraint(s) {
	case ConstraintNone:
		return ConstraintNone, nil
	case ConstraintFriendsOnly:
		return ConstraintFriendsOnly, nil
	case ConstraintInvitationOnly:
		return ConstraintInvitationOnly, nil
	default:
		return ConstraintNone, errors.New("invalid constraint")

	}
}

type Session struct {
	id         string
	users      map[string]*user.User
	capacity   int
	minRating  int
	maxRating  int
	constraint Constraint
	owner      *user.User
	messages   []*message.Message
	movements  []*movement.Movement
}

func New(owner *user.User, capacity int, minRating int, maxRating int, constraint Constraint) *Session {
	return &Session{
		id:         shortuuid.New(),
		users:      map[string]*user.User{},
		capacity:   capacity,
		minRating:  minRating,
		maxRating:  maxRating,
		constraint: constraint,
		owner:      owner,
		messages:   nil,
		movements:  nil,
	}
}

func (s *Session) SetUsers(users map[string]*user.User) {
	s.users = users
}

func (s *Session) SetCapacity(capacity int) {
	s.capacity = capacity
}

func (s *Session) SetMinRating(minRating int) {
	s.minRating = minRating
}

func (s *Session) MaxRating() int {
	return s.maxRating
}

func (s *Session) SetMaxRating(maxRating int) {
	s.maxRating = maxRating
}

func (s *Session) SetConstraint(constraint Constraint) {
	s.constraint = constraint
}

func (s *Session) SetOwner(owner *user.User) {
	s.owner = owner
}

func (s *Session) Users() map[string]*user.User {
	return s.users
}

func (s *Session) Capacity() int {
	return s.capacity
}

func (s *Session) Id() string {
	return s.id
}

func (s *Session) MinRating() int {
	return s.minRating
}

func (s *Session) Constraint() string {
	return string(s.constraint)
}

func (s *Session) Owner() *user.User {
	return s.owner
}

func (s *Session) AddUser(u *user.User) error {
	_, exists := s.users[u.Id()]
	if exists {
		return errors.New("user already exists")
	}

	s.users[u.Id()] = u

	return nil
}

func (s *Session) RemoveUser(userId string) error {
	_, exists := s.users[userId]
	if !exists {
		return errors.New("user does not exist")
	}

	delete(s.users, userId)
	for _, u := range s.users {
		s.owner = u
		break
	}

	return nil
}

func (s *Session) Messages() []*message.Message {
	return s.messages
}

func (s *Session) AddMessage(message *message.Message) error {
	s.messages = append(s.messages, message)
	return nil
}

func (s *Session) Movements() []*movement.Movement {
	return s.movements
}

func (s *Session) AddMovement(movement *movement.Movement) error {
	s.movements = append(s.movements, movement)
	return nil
}
