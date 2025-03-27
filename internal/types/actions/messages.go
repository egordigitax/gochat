package actions

import (
	"chat-service/internal/types"
)

type RequestMessageAction struct {
	ChatUid   string `json:"chat_uid"`
	AuthorUid string `json:"author_uid"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

func InitRequestMessageAction(msg types.Message) RequestMessageAction {
	return RequestMessageAction{
		ChatUid:   msg.ChatUid,
		AuthorUid: msg.UserUid,
		CreatedAt: msg.CreatedAt,
		Text:      msg.Text,
	}
}

func (g RequestMessageAction) GetActionType() types.ActionType {
	return types.REQUEST_MESSAGE
}

type SendMessageAction struct {
	ChatUid   string `json:"chat_uid"`
	AuthorUid string `json:"author_uid"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

func (s SendMessageAction) GetActionType() types.ActionType {
	return types.SEND_MESSAGE
}
