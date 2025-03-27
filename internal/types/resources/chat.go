package resources

import (
	"chat-service/internal/types"
)

type Chat struct {
	Uid         string  `json:"uid"`
	Title       string  `json:"title"`
	MediaUrl    string  `json:"media_url"`
	UnreadCount int     `json:"unread_count"`
	LastMessage Message `json:"message"`
	UpdatedAt   string  `json:"updated_at"`
	Status      string  `json:"status"`
}

func NewChatFromEntity(entity types.Chat) Chat {
	return Chat{
		Uid:       entity.Uid,
		Title:     entity.Title,
		UpdatedAt: entity.UpdatedAt,
		LastMessage: Message{
			Username:  entity.LastMessage.UserInfo.Nickname,
			AuthorUid: entity.LastMessage.UserInfo.Uid,
			ChatUid:   entity.LastMessage.ChatUid,
			Text:      entity.LastMessage.Text,
			CreatedAt: entity.LastMessage.CreatedAt,
		},
	}
}
