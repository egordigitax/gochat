package repositories

import "chat-service/internal/domain/entities"

type ChatsService interface {
	CheckIfUserHasAccess(
		user_uid string,
		chat_uid string,
	) (bool, error)
}

type ChatsStorage interface {
	GetChatsByUserUid(
		userUid string,
		limit, offset int,
	) ([]entities.Chat, error)
	GetChatByUid(
		chat_uid string,
	) (entities.Chat, error)
	CheckIfUserHasAccess(
		user_uid string,
		chat_uid string,
	) (bool, error)
	FetchChatsLastMessages(
		chats *[]entities.Chat,
	) error
	GetAllUsersFromChatByUid(
		chat_uid string,
	) ([]string, error)
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
