package chat_list

import (
	"chat-service/internal/application/common/ports"
	"chat-service/internal/application/schema/dto"
	"chat-service/internal/application/schema/resources"
	"github.com/spf13/viper"
	"log"
)

type ChatsClient struct {
	Hub    *ChatsHub
	Conn   ports.ClientTransport
	UserId string
	Send   chan resources.BaseMessage
}

func NewChatsClient(hub *ChatsHub, conn ports.ClientTransport, userId string) *ChatsClient {
	sendChan := make(
		chan resources.BaseMessage,
		viper.GetInt("app.users_msg_buff"),
	)

	return &ChatsClient{
		Hub:    hub,
		Conn:   conn,
		UserId: userId,
		Send:   sendChan,
	}
}

func (c *ChatsClient) GetChats() {
	chats, err := c.Hub.chats.GetChatsByUserUid(c.UserId)
	if err != nil {
		log.Println(err)
		return
	}

	data := dto.GetUserChatsPayload{
		Items: chats,
	}

	c.Send <- resources.BaseMessage{
		Action: data.GetActionType(),
		Data:   data,
	}
}
