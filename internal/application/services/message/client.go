package message

import (
	"chat-service/internal/application/constants"
	"chat-service/internal/application/ports"
	"chat-service/internal/domain/entities"
	"chat-service/internal/schema/dto"
	"context"
	"log"
	"slices"

	"github.com/spf13/viper"
)

type MessagesClient struct {
	Hub     *MessagesHub
	Conn    ports.ClientTransport
	UserUid string
	ChatUid string
	Send    chan dto.SendMessageToClientPayload
}

func NewMessagesClient(
	hub *MessagesHub,
	Conn ports.ClientTransport,
	UserUid string, ChatUid string,
) *MessagesClient {

	sendChan := make(
		chan dto.SendMessageToClientPayload,
		viper.GetInt("app.users_msg_buff"),
	)

	return &MessagesClient{
		Hub:     hub,
		Conn:    Conn,
		UserUid: UserUid,
		ChatUid: ChatUid,
		Send:    sendChan,
	}
}

func (c *MessagesClient) GetMessageFromClient(
	ctx context.Context,
	msg dto.GetMessageFromClientPayload,
) {
	message := msg.ToEntity()

	err := c.Hub.broker.SendMessageToQueue(
		ctx,
		constants.CHATS_QUEUE,
		message,
	)
	if err != nil {
		log.Println(err.Error())
	}
}

func (c *MessagesClient) SendMessageToClient(
	ctx context.Context,
	msg entities.Message,
) error {
	c.Send <- dto.BuildSendMessageToClientPayloadFromEntity(msg)

	return nil
}

func (c *MessagesClient) SendMessagesHistory(limit, offset int) {
	history, err := c.Hub.messages.GetMessagesHistory(c.ChatUid, limit, offset)
	if err != nil {
		log.Println("smth wrong:", err)
		return
	}

	// TODO: ordering or handle it in GetMessagesHistory

	slices.Reverse(history)

	for _, msg := range history {
		c.Send <- dto.BuildSendMessageToClientPayloadFromEntity(msg)
	}
}

func (c *MessagesClient) GetMe() entities.User {
	// return c.H
	panic("implement me")
}
