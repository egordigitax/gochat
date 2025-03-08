package repositories

import "chat-service/internal/domain/entities"


type ChatsService interface {
	CheckIfUserHasAccess(
		user_uid string,
		chat_uid string,
	) (bool, error)
}

type ChatsStorage interface {
	GetUsersChats(
		user_uid string,
		limit int,
		offset int,
	) ([]entities.Chat, error)
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
	) ([]entities.Chat, error)
	SetUsersChats(
		user_uid string,
		chats []entities.Chat,
	) error
}
