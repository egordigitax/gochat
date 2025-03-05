package interfaces

import "chat-service/internal/domain"

type ChatsStorage interface {
	GetUsersChats(
		user_uid string,
		limit int,
		offset int,
	) ([]domain.Chat, error)
}

type ChatsCache interface {
	GetUsersChats(
		user_uid string,
		limit int,
		offset int,
	) ([]domain.Chat, error)
}
