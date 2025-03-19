package dto

import (
	"chat-service/internal/application/schema/resources"
)

type RequestMessagePayload struct {
	ChatUid   string `json:"chat_uid"`
	AuthorUid string `json:"author_uid"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

func (g RequestMessagePayload) GetActionType() resources.ActionType {
	return resources.REQUEST_MESSAGE
}

type SendMessagePayload struct {
	ChatUid   string `json:"chat_uid"`
	AuthorUid string `json:"author_uid"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

func (s SendMessagePayload) GetActionType() resources.ActionType {
	return resources.SEND_MESSAGE
}
