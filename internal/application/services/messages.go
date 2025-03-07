package services

import (
	"chat-service/internal/domain"
	"chat-service/internal/domain/interfaces"
)

type MessageService struct {
	MessagesStorage interfaces.MessagesStorage
	MessagesCache   interfaces.MessagesCache
	ChatsStorage    interfaces.ChatsStorage
	ChatsCache      interfaces.ChatsCache
}

func NewMessageService(
	messagesStorage interfaces.MessagesStorage,
	messagesCache interfaces.MessagesCache,
	chatsStorage interfaces.ChatsStorage,
	chatsCache interfaces.ChatsCache,
) *MessageService {
	return &MessageService{
		MessagesStorage: messagesStorage,
		MessagesCache:   messagesCache,
		ChatsStorage:    chatsStorage,
		ChatsCache:      chatsCache,
	}
}

func (m *MessageService) GetChatsList(user_uid string) ([]domain.Chat, error) {
	msgs, err := m.ChatsStorage.GetUsersChats(user_uid, 10, 0)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (m *MessageService) CreateNewChat(
	title string,
	chat_type string,
	media_url string,
	users_uids []string,
) {
}
