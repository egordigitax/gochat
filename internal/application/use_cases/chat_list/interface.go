package chat_list

import "chat-service/internal/schema/dto"

type IChatsService interface {
	CheckIfUserHasAccess(userUid string, chatUid string) (bool, error)
	GetChatsByUserUid(payload dto.GetUserChatsByUidPayload) (dto.GetUserChatsByUidResponse, error)
	GetAllUsersFromChatByUid(chatUid string) ([]string, error)
	CreateNewChat(title string, chatType string, mediaUrl string, usersUids []string)
}
