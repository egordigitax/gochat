package dto

import (
	"chat-service/internal/application/schema/resources"
	"chat-service/internal/domain/entities"
)

type RequestUserChatsPayload struct {
	Items []resources.Chat `json:"items"`
}

func BuildRequestUserChatsPayloadFromResources(
	chats []resources.Chat,
) RequestUserChatsPayload {
	return RequestUserChatsPayload{
		Items: chats,
	}
}

func BuildRequestUserChatsPayloadFromEntities(
	chats []entities.Chat,
) RequestUserChatsPayload {

    resChat := make([]resources.Chat, len(chats))
    for i, item := range chats {
        resChat[i] = resources.NewChatFromEntity(item)
    }

	return RequestUserChatsPayload{
		Items: resChat,
	}
}

func (s RequestUserChatsPayload) GetActionType() resources.ActionType {
	return resources.REQUEST_CHATS
}
