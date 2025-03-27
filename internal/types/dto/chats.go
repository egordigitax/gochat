package dto

import (
	"chat-service/internal/types"
	resources2 "chat-service/internal/types/resources"
)

type RequestUserChatsPayload struct {
	Items []resources2.Chat `json:"items"`
}

func BuildRequestUserChatsPayloadFromResources(
	chats []resources2.Chat,
) RequestUserChatsPayload {
	return RequestUserChatsPayload{
		Items: chats,
	}
}

func BuildRequestUserChatsPayloadFromEntities(
	chats []types.Chat,
) RequestUserChatsPayload {

	resChat := make([]resources2.Chat, len(chats))
	for i, item := range chats {
		resChat[i] = resources2.NewChatFromEntity(item)
	}

	return RequestUserChatsPayload{
		Items: resChat,
	}
}

func (s RequestUserChatsPayload) GetActionType() types.ActionType {
	return types.REQUEST_CHATS
}
