package chat_list

import (
	"chat-service/internal/types"
)

type IChatsService interface {
	CheckIfUserHasAccess(userUid string, chatUid string) (bool, error)
	GetChatsByUserUid(userUid string) ([]types.Chat, error)
	GetUsersFromChat(chatUid string) ([]string, error)
	CreateNewChat(title string, chatType string, mediaUrl string, usersUids []string)
}
