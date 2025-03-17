package message

import "chat-service/internal/domain/entities"

type IMessagesService interface {
	GetMessagesHistory(chatUID string, limit, offset int) ([]entities.Message, error)
	SaveMessagesBulk(msgs ...entities.Message) error
}
