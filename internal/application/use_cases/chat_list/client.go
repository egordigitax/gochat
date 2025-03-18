package chat_list

import (
	"chat-service/internal/application/common/ports"
	"chat-service/internal/schema/dto"
	"chat-service/internal/schema/resources"
	"github.com/spf13/viper"
	"log"
)

type ChatsClient struct {
	Hub    *ChatsHub
	Conn   ports.ClientTransport
	UserID string
	Send   chan resources.BaseMessage
}

func NewChatsClient(hub *ChatsHub, conn ports.ClientTransport, userID string) *ChatsClient {
	sendChan := make(
		chan resources.BaseMessage,
		viper.GetInt("app.users_msg_buff"),
	)

	return &ChatsClient{
		Hub:    hub,
		Conn:   conn,
		UserID: userID,
		Send:   sendChan,
	}
}

func (c *ChatsClient) GetChats() {
	chats, err := c.Hub.chats.GetChatsByUserUid(
		dto.GetUserChatsByUidPayload{
			UserUid: c.UserID,
		})
	if err != nil {
		log.Println(err)
		return
	}
	c.Send <- resources.BaseMessage{
		Action: "get_chats",
		Data:   chats,
	}
}

func (c *ChatsClient) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
