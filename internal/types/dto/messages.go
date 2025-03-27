package dto

import (
	"chat-service/internal/types"
	resources2 "chat-service/internal/types/resources"
)

type RequestMessagePayload struct {
	ChatUid   string `json:"chat_uid"`
	AuthorUid string `json:"author_uid"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

func BuildRequestMessagePayloadFromEntity(msg types.Message) RequestMessagePayload {
	return RequestMessagePayload{
		ChatUid:   msg.ChatUid,
		AuthorUid: msg.UserUid,
		CreatedAt: msg.CreatedAt,
		Text:      msg.Text,
	}
}

func BuildRequestMessagePayloadFromResources(msg resources2.Message) RequestMessagePayload {
	return RequestMessagePayload{
		ChatUid:   msg.ChatUid,
		AuthorUid: msg.AuthorUid,
		CreatedAt: msg.CreatedAt,
		Text:      msg.Text,
	}
}

func (g RequestMessagePayload) GetActionType() types.ActionType {
	return types.REQUEST_MESSAGE
}

type SendMessagePayload struct {
	ChatUid   string `json:"chat_uid"`
	AuthorUid string `json:"author_uid"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

func (s SendMessagePayload) GetActionType() types.ActionType {
	return types.SEND_MESSAGE
}
