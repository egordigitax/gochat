package dto

import "chat-service/internal/schema/resources"

type GetUserChatsByUidResponse struct {
	Items []resources.Chat `json:"items"`
}

type GetUserChatsByUidPayload struct {
	UserUid string `json:"user_uid"`
}
