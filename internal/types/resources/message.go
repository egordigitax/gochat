package resources

import (
	"chat-service/internal/types"
)

type Message struct {
	Username  string `json:"username"`
	AuthorUid string `json:"author_uid"`
	ChatUid   string `json:"chat_uid"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

func NewMessageFromEntity(msg types.Message) Message {
	return Message{
		Username:  msg.UserInfo.Nickname,
		AuthorUid: msg.UserUid,
		ChatUid:   msg.ChatUid,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
}
