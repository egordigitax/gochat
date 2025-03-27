package chat_list

import (
	"chat-service/internal/memory"
	"chat-service/internal/storage"
	"chat-service/internal/types"
)

type ChatsService struct {
	ChatsStorage storage.ChatsStorage
	ChatsCache   memory.ChatsCache
}

func NewChatsService(
	chatsStorage storage.ChatsStorage,
	chatsCache memory.ChatsCache,
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
) ([]types.Chat, error) {

	chats, err := m.ChatsStorage.GetChatsByUserUid(userUid, 10, 0)
	if err != nil {
		return chats, err
	}

	err = m.ChatsStorage.FetchChatsLastMessages(&chats)
	if err != nil {
		return chats, err
	}

	return chats, nil
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
