package dto

import (
	"chat-service/internal/application/schema/resources"
	"chat-service/internal/domain/entities"
)

type RequestMessagePayload struct {
	ChatUid   string `json:"chat_uid"`
	AuthorUid string `json:"author_uid"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

func BuildRequestMessagePayloadFromEntity(msg entities.Message) RequestMessagePayload {
    return RequestMessagePayload{
    	ChatUid:   msg.ChatUid,
    	AuthorUid: msg.UserUid,
    	CreatedAt: msg.CreatedAt,
    	Text:      msg.Text,
    }
}

func BuildRequestMessagePayloadFromResources(msg resources.Message) RequestMessagePayload {
    return RequestMessagePayload{
    	ChatUid:   msg.ChatUid,
    	AuthorUid: msg.AuthorUid,
    	CreatedAt: msg.CreatedAt,
    	Text:      msg.Text,
    }
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
