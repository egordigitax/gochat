package chat_list

import (
	"chat-service/internal/application/schema/resources"
)

type IChatsService interface {
	CheckIfUserHasAccess(userUid string, chatUid string) (bool, error)
	GetChatsByUserUid(userUid string) ([]resources.Chat, error)
	GetUsersFromChat(chatUid string) ([]string, error)
	CreateNewChat(title string, chatType string, mediaUrl string, usersUids []string)
}
