package domain

type ChatRepository interface {
	SaveMessage(msg Message) error
	GetMessages(chatID string, limit int) ([]Message, error)
}
