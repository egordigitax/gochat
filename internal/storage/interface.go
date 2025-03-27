package storage

import "chat-service/internal/types"

type ChatsStorage interface {
	GetChatsByUserUid(
		userUid string,
		limit, offset int,
	) ([]types.Chat, error)
	GetChatByUid(
		chat_uid string,
	) (types.Chat, error)
	CheckIfUserHasAccess(
		user_uid string,
		chat_uid string,
	) (bool, error)
	FetchChatsLastMessages(
		chats *[]types.Chat,
	) error
	GetAllUsersFromChatByUid(
		chat_uid string,
	) ([]string, error)
}

type UsersStorage interface {
	GetUserByUid(userUid string) (types.User, error)
}
