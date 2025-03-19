package messages

import (
	"chat-service/internal/application/common/constants"
	"chat-service/internal/application/common/ports"
	"chat-service/internal/application/schema/dto"
	"chat-service/internal/application/schema/resources"
	"chat-service/internal/domain/entities"
	"context"
	"log"

	"github.com/spf13/viper"
)

type MessageClient struct {
	Hub     *MessageHub
	Conn    ports.ClientTransport
	UserUid string
	ChatUid string
	Send    chan resources.Action
}

func NewMessagesClient(
	hub *MessageHub,
	Conn ports.ClientTransport,
	UserUid string, ChatUid string,
) *MessageClient {

	sendChan := make(
		chan resources.Action,
		viper.GetInt("app.users_msg_buff"),
	)

	return &MessageClient{
		Hub:     hub,
		Conn:    Conn,
		UserUid: UserUid,
		ChatUid: ChatUid,
		Send:    sendChan,
	}
}

func (c *MessageClient) SendMessage(
	ctx context.Context,
	msg dto.SendMessagePayload,
) {

	message := entities.NewMessage(
		msg.Text,
		msg.AuthorUid,
		msg.ChatUid,
	)

	err := c.Hub.broker.SendMessageToQueue(
		ctx,
		constants.CHATS_QUEUE,
		message,
	)
	if err != nil {
		log.Println(err)
	}
}

func (c *MessageClient) RequestMessage(
	msg resources.Message,
) error {

	data := dto.BuildRequestMessagePayloadFromResources(msg)
	c.Send <- resources.BuildAction(data)

	return nil
}

func (c *MessageClient) RequestMessageHistory(limit, offset int) {
	history, err := c.Hub.messages.GetMessagesHistory(c.ChatUid, limit, offset)
	if err != nil {
		log.Println("error while fetching history:", err)
		return
	}

	for _, msg := range history {
		data := dto.BuildRequestMessagePayloadFromEntity(msg)
		c.Send <- resources.BuildAction(data)
	}
}

func (c *MessageClient) GetMe() resources.User {
	panic("implement me")
}
