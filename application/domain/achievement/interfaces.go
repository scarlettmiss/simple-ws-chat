package achievement

type Service interface {
	Achievement(id string) (*Achievement, error)
	CreateAchievement(name string) (*Achievement, error)
}

type Repository interface {
	CreateAchievement(name string) (*Achievement, error)
	Achievement(id string) (*Achievement, error)
	Achievements() map[string]*Achievement
	AchievementByName(name string) (*Achievement, error)
}
