package dto

import (
	"chat-service/internal/application/schema/resources"
)

type GetUserChatsPayload struct {
	Items []resources.Chat `json:"items"`
}

func (s GetUserChatsPayload) GetActionType() resources.ActionType {
	return "GetChats"
}
