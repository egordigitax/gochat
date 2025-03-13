package services

import (
	"chat-service/internal/domain/repositories"
	"chat-service/internal/schema/dto"
	"chat-service/internal/schema/resources"
)

type ChatsService struct {
	ChatsStorage repositories.ChatsStorage
	ChatsCache   repositories.ChatsCache
}

func NewChatsService(
	chatsStorage repositories.ChatsStorage,
	chatsCache repositories.ChatsCache,
) *ChatsService {
	return &ChatsService{
		ChatsStorage: chatsStorage,
		ChatsCache:   chatsCache,
	}
}

func (m *ChatsService) CheckIfUserHasAccess(user_uid string, chat_uid string) (bool, error) {
	return m.ChatsStorage.CheckIfUserHasAccess(user_uid, chat_uid)
}

func (m *ChatsService) GetChatsByUserUid(
	payload dto.GetUserChatsByUidPayload,
) (dto.GetUserChatsByUidResponse, error) {

	response := dto.GetUserChatsByUidResponse{}

	chats, err := m.ChatsStorage.GetChatsByUserUid(payload.UserUid, 10, 0)
	if err != nil {
		return response, err
	}

	err = m.ChatsStorage.FetchChatsLastMessages(&chats)
	if err != nil {
		return response, err
	}

	response.Items = make([]resources.Chat, len(chats))

	for i, item := range chats {
		response.Items[i] = resources.Chat{}
		response.Items[i].FromEnitity(&item)
	}

	return response, nil
}

func (m *ChatsService) GetAllUsersFromChatByUid(chat_uid string) ([]string, error) {
    // get from 
	return m.ChatsStorage.GetAllUsersFromChatByUid(chat_uid)
}

func (m *ChatsService) CreateNewChat(
	title string,
	chat_type string,
	media_url string,
	users_uids []string,
) {
}
