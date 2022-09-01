package movement

import (
	"github.com/lithammer/shortuuid"
	"time"
)

type Movement struct {
	id           string
	movementType string
	userId       string
	acceleration *Point
	position     *Point
	createdOn    time.Time
}

type Point struct {
	X int
	Y int
}

func New(movementType string, userId string, acceleration *Point, position *Point) *Movement {
	timestamp := time.Now()
	return &Movement{
		id:           shortuuid.New(),
		movementType: movementType,
		userId:       userId,
		acceleration: acceleration,
		position:     position,
		createdOn:    timestamp,
	}
}

func (m *Movement) Id() string {
	return m.id
}

func (m *Movement) MovementType() string {
	return m.movementType
}

func (m *Movement) SetMovementType(movementType string) {
	m.movementType = movementType
}

func (m *Movement) UserId() string {
	return m.userId
}

func (m *Movement) SetUserId(userId string) {
	m.userId = userId
}

func (m *Movement) Acceleration() *Point {
	return m.acceleration
}

func (m *Movement) SetAcceleration(acceleration *Point) {
	m.acceleration = acceleration
}

func (m *Movement) Position() *Point {
	return m.position
}

func (m *Movement) SetPosition(position *Point) {
	m.position = position
}

func (m *Movement) CreatedOn() time.Time {
	return m.createdOn
}
