package interfaces

import "chat-service/internal/domain"

type MessageService interface {

}

type MessagesStorage interface {
	GetMessages(
		chat_uid string,
		limit, offset int,
	) ([]domain.Message, error)

	SaveMessage(msg domain.Message) error
}

type MessagesCache interface {
	GetMessages(
		chat_uid string,
		limit, offset int,
	) ([]domain.Message, error)

	SaveMessage(msg domain.Message) error
}
