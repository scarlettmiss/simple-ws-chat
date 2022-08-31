package message

type Service interface {
	CreateMessage(userId string, message string) (*Message, error)
	Messages() map[string]*Message
	MessagesByUserId(userId string) map[string]*Message
}

type Repository interface {
	CreateMessage(userId string, message string) (*Message, error)
	Message(id string) (*Message, error)
	Messages() map[string]*Message
	MessagesByUserId(userId string) map[string]*Message
}
