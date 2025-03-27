package memory

import "chat-service/internal/types"

type ChatsCache interface {
	GetUsersChats(
		user_uid string,
		limit int,
		offset int,
	) ([]types.Chat, error)
	SetUsersChats(
		user_uid string,
		chats []types.Chat,
	) error
}
