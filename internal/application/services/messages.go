package services

import "chat-service/internal/domain/repositories"

type MessageService struct {
	MessagesStorage repositories.MessagesStorage
	MessagesCache   repositories.MessagesCache
}

func NewMessageService(
	messagesStorage repositories.MessagesStorage,
	messagesCache repositories.MessagesCache,
) *MessageService {
	return &MessageService{
		MessagesStorage: messagesStorage,
		MessagesCache:   messagesCache,
	}
}

var _ repositories.MessageService = MessageService{}
