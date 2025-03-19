package dto

import (
	"chat-service/internal/application/schema/resources"
)

type RequestUserChatsPayload struct {
	Items []resources.Chat `json:"items"`
}

func (s RequestUserChatsPayload) GetActionType() resources.ActionType {
	return resources.REQUEST_CHATS
}
