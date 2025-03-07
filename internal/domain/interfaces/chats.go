package interfaces

import "chat-service/internal/domain"

type ChatsStorage interface {
	GetUsersChats(
		user_uid string,
		limit int,
		offset int,
	) ([]domain.Chat, error)
	CheckIfUserHasAccess(
		user_uid string,
		chat_uid string,
	) (bool, error)
}

type ChatsCache interface {
	GetUsersChats(
		user_uid string,
		limit int,
		offset int,
	) ([]domain.Chat, error)
	SetUsersChats(
		user_uid string,
		chats []domain.Chat,
	) error
}
