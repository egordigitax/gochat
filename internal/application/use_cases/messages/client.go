package messages

import (
	"chat-service/internal/application/common/constants"
	"chat-service/internal/application/common/ports"
	"chat-service/internal/application/schema/dto"
	"chat-service/internal/application/schema/resources"
	"chat-service/internal/domain/entities"
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
	Send    chan resources.BaseMessage
}

func NewMessagesClient(
	hub *MessageHub,
	Conn ports.ClientTransport,
	UserUid string, ChatUid string,
) *MessageClient {

	sendChan := make(
		chan resources.BaseMessage,
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

func (c *MessageClient) GetMessage(
	msg resources.Message,
) error {

	data := dto.GetMessagePayload{
		AuthorUid: msg.AuthorUid,
		ChatUid:   msg.ChatUid,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}

	c.Send <- resources.BaseMessage{
		Action: data.GetActionType(),
		Data:   data,
	}

	return nil
}

func (c *MessageClient) GetMessagesHistory(limit, offset int) {
	history, err := c.Hub.messages.GetMessagesHistory(c.ChatUid, limit, offset)
	if err != nil {
		log.Println("smth wrong:", err)
		return
	}

	// TODO: ordering or handle it in GetMessagesHistory

	slices.Reverse(history)

	for _, msg := range history {
		data := dto.GetMessagePayload{
			AuthorUid: msg.UserUid,
			ChatUid:   msg.ChatUid,
			Text:      msg.Text,
			CreatedAt: msg.CreatedAt,
		}

		c.Send <- resources.BaseMessage{
			Action: data.GetActionType(),
			Data:   data,
		}
	}
}

func (c *MessageClient) GetMe() resources.User {
	// return c.H
	panic("implement me")
}
