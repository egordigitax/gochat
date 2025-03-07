package services

import (
	"chat-service/internal/domain"
	"chat-service/internal/domain/interfaces"
)

type ChatsService struct {
	ChatsStorage interfaces.ChatsStorage
	ChatsCache   interfaces.ChatsCache
}

func NewChatsService(
	chatsStorage interfaces.ChatsStorage,
	chatsCache interfaces.ChatsCache,
) *ChatsService {
	return &ChatsService{
		ChatsStorage: chatsStorage,
		ChatsCache:   chatsCache,
	}
}

func (m ChatsService) CheckIfUserHasAccess(user_uid string, chat_uid string) (bool, error) {
	return m.ChatsStorage.CheckIfUserHasAccess(user_uid, chat_uid)
}

func (m *ChatsService) GetChatsList(user_uid string) ([]domain.Chat, error) {
	msgs, err := m.ChatsStorage.GetUsersChats(user_uid, 10, 0)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (m *ChatsService) CreateNewChat(
	title string,
	chat_type string,
	media_url string,
	users_uids []string,
) {
}
