package chat_list

import (
	"chat-service/internal/application/common/ports"
	"chat-service/internal/schema/dto"
	"log"
)

type ChatsClient struct {
	Hub    *ChatsHub
	Conn   ports.ClientTransport
	UserID string
	Send   chan dto.GetUserChatsByUidResponse
}

func (c *ChatsClient) UpdateChats() {
	chats, err := c.Hub.chats.GetChatsByUserUid(
		dto.GetUserChatsByUidPayload{
			UserUid: c.UserID,
		})
	if err != nil {
		log.Println(err)
		return
	}
	c.Send <- chats
}

func (c *ChatsClient) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
