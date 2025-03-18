package messages

import "chat-service/internal/domain/entities"

type IMessageService interface {
	GetMessagesHistory(chatUID string, limit, offset int) ([]entities.Message, error)
	SaveMessagesBulk(msgs ...entities.Message) error
}
