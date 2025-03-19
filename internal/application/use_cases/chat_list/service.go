package chat_list

import (
	"chat-service/internal/application/schema/resources"
	"chat-service/internal/domain/repositories"
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

func (m *ChatsService) CheckIfUserHasAccess(userUid string, chatUid string) (bool, error) {
	return m.ChatsStorage.CheckIfUserHasAccess(userUid, chatUid)
}

func (m *ChatsService) GetChatsByUserUid(
	userUid string,
) ([]resources.Chat, error) {

	var response []resources.Chat

	chats, err := m.ChatsStorage.GetChatsByUserUid(userUid, 10, 0)
	if err != nil {
		return response, err
	}

	err = m.ChatsStorage.FetchChatsLastMessages(&chats)
	if err != nil {
		return response, err
	}

	response = make([]resources.Chat, len(chats))

	for i, item := range chats {
		response[i] = resources.NewChatFromEntity(item)
	}

	return response, nil
}

func (m *ChatsService) GetUsersFromChat(chatUid string) ([]string, error) {
	return m.ChatsStorage.GetAllUsersFromChatByUid(chatUid)
}

func (m *ChatsService) CreateNewChat(
	title string,
	chat_type string,
	media_url string,
	users_uids []string,
) {
	panic("implement me")
}
