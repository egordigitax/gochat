package actions

import (
	"chat-service/internal/types"
)

type RequestUserChatsAction struct {
	Items []types.Chat `json:"items"`
}

func (s RequestUserChatsAction) GetActionType() types.ActionType {
	return types.REQUEST_CHATS
}

func InitRequestUserChatsAction(chats []types.Chat) RequestUserChatsAction {
	return RequestUserChatsAction{Items: chats}
}
