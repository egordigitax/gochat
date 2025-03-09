package resources

import "chat-service/internal/domain/entities"

type Chat struct {
	Title       string  `json:"title"`
	MediaUrl    string  `json:"media_url"`
	UnreadCount int     `json:"unread_count"`
	LastMessage Message `json:"message"`
	UpdatedAt   string  `json:"updated_at"`
	Status      string  `json:"status"`
}

func (c *Chat) FromEnitity(entity *entities.Chat) {
	c.Title = entity.Title
	c.UpdatedAt = entity.UpdatedAt
	c.LastMessage = Message{
		Username:  entity.LastMessage.UserInfo.Nickname,
		AuthorUid: entity.LastMessage.UserInfo.Uid,
		ChatUid:   entity.Uid,
		Text:      entity.LastMessage.Text,
		CreatedAt: entity.LastMessage.CreatedAt,
	}
}
