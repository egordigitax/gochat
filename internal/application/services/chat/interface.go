package chat

import "chat-service/internal/schema/dto"

type IChatsService interface {
	CheckIfUserHasAccess(user_uid string, chat_uid string) (bool, error)
	GetChatsByUserUid(payload dto.GetUserChatsByUidPayload) (dto.GetUserChatsByUidResponse, error)
	GetAllUsersFromChatByUid(chat_uid string) ([]string, error)
	CreateNewChat(title string, chat_type string, media_url string, users_uids []string)
}
