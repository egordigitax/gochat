package messages

import (
	"chat-service/internal/types"
)

type IMessageService interface {
	GetMessagesHistory(chatUid string, limit, offset int) ([]types.Message, error)
	SaveMessagesBulk(msgs ...types.Message) error
}

type MessagesStorage interface {
	GetMessages(
		chat_uid string,
		limit, offset int,
	) ([]types.Message, error)

	SaveMessage(msg types.Message) error
	SaveMessagesBulk(msg ...types.Message) error
}

type MessagesCache interface {
	GetMessagesByChatUid(
		chat_uid string,
	) ([]types.Message, error)

	DeleteMessage(msg types.Message) error
	SaveMessage(msg types.Message) error
}
