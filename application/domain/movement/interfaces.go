package movement

type Service interface {
	CreateMovement(movementType string, userId string, acceleration *Point, position *Point) (*Movement, error)
	Movements() map[string]*Movement
	MovementsByUserId(userId string) map[string]*Movement
}

type Repository interface {
	CreateMovement(movementType string, userId string, acceleration *Point, position *Point) (*Movement, error)
	Movement(id string) (*Movement, error)
	Movements() map[string]*Movement
	MovementsByUserId(userId string) map[string]*Movement
}
