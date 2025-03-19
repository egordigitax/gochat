package dto

import (
	"chat-service/internal/application/schema/resources"
)

type SendMessagePayload struct {
	ChatUid   string `json:"chat_uid"`
	AuthorUid string `json:"author_uid"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

func (g SendMessagePayload) GetActionType() resources.ActionType {
	return "GetMessage"
}

type GetMessagePayload struct {
	ChatUid   string `json:"chat_uid"`
	AuthorUid string `json:"author_uid"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

func (s GetMessagePayload) GetActionType() resources.ActionType {
	return "SendMessage"
}
