package storage

import "chat-service/internal/types"

type ChatsStorage interface {
	GetChatsByUserUid(
		userUid string,
		limit, offset int,
	) ([]types.Chat, error)
	GetChatByUid(
		chatUid string,
	) (types.Chat, error)
	CheckIfUserHasAccess(
		userUid string,
		chatUid string,
	) (bool, error)
	FetchChatsLastMessages(
		chats *[]types.Chat,
	) error
	GetAllUsersFromChatByUid(
		chatUid string,
	) ([]string, error)
}

type UsersStorage interface {
	GetUserByUid(userUid string) (types.User, error)
}
