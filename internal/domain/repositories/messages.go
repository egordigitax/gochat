package repositories

import "chat-service/internal/domain/entities"

type MessageService interface {
}

type MessagesStorage interface {
	GetMessages(
		chat_uid string,
		limit, offset int,
	) ([]entities.Message, error)

	SaveMessage(msg entities.Message) error
	SaveMessagesBulk(msg ...entities.Message) error
}

type MessagesCache interface {
	GetMessagesByChatUid(
		chat_uid string,
	) ([]entities.Message, error)

	DeleteMessage(msg entities.Message) error
	SaveMessage(msg entities.Message) error
}
