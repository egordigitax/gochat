package chat_list

import (
	"chat-service/common/ports"
	"chat-service/internal/types"
	"chat-service/internal/types/dto"
	"github.com/spf13/viper"
	"log"
)

type ChatsClient struct {
	Hub    *ChatsHub
	Conn   ports.ClientTransport
	UserId string
	Send   chan types.Action
}

func NewChatsClient(hub *ChatsHub, conn ports.ClientTransport, userId string) *ChatsClient {
	sendChan := make(
		chan types.Action,
		viper.GetInt("app.users_msg_buff"),
	)

	return &ChatsClient{
		Hub:    hub,
		Conn:   conn,
		UserId: userId,
		Send:   sendChan,
	}
}

func (c *ChatsClient) RequestChats() {
	chats, err := c.Hub.chats.GetChatsByUserUid(c.UserId)
	if err != nil {
		log.Println(err)
		return
	}

	data := dto.BuildRequestUserChatsPayloadFromEntities(chats)
	c.Send <- types.BuildAction(data)
}
