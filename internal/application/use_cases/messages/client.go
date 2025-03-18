package messages

import (
	"chat-service/internal/application/common/constants"
	"chat-service/internal/application/common/ports"
	"chat-service/internal/domain/entities"
	"chat-service/internal/schema/dto"
	"context"
	"log"
	"slices"

	"github.com/spf13/viper"
)

type MessageClient struct {
	Hub     *MessageHub
	Conn    ports.ClientTransport
	UserUid string
	ChatUid string
	Send    chan dto.SendMessageToClientPayload
}

func NewMessagesClient(
	hub *MessageHub,
	Conn ports.ClientTransport,
	UserUid string, ChatUid string,
) *MessageClient {

	sendChan := make(
		chan dto.SendMessageToClientPayload,
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

func (c *MessageClient) GetMessageFromClient(
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

func (c *MessageClient) SendMessageToClient(
	ctx context.Context,
	msg entities.Message,
) error {
	c.Send <- dto.BuildSendMessageToClientPayloadFromEntity(msg)

	return nil
}

func (c *MessageClient) SendMessagesHistory(limit, offset int) {
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

func (c *MessageClient) GetMe() entities.User {
	// return c.H
	panic("implement me")
}
