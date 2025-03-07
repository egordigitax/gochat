package services

import (
	"chat-service/internal/domain/interfaces"
)

type MessageService struct {
	MessagesStorage interfaces.MessagesStorage
	MessagesCache   interfaces.MessagesCache
}

func NewMessageService(
	messagesStorage interfaces.MessagesStorage,
	messagesCache interfaces.MessagesCache,
) *MessageService {
	return &MessageService{
		MessagesStorage: messagesStorage,
		MessagesCache:   messagesCache,
	}
}

var _ interfaces.MessageService = MessageService{}
