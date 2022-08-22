package achievment

type Service interface {
	Achivement(id string) (*Achievement, error)
	CreateAchivement(username string, password string) (*Achievement, error)
}

type Repository interface {
	CreateAchivement(name string) (*Achievement, error)
	Achivement(id string) (*Achievement, error)
	Achivements() map[string]*Achievement
}
