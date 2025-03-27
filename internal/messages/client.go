package messages

import (
	"chat-service/common/constants"
	"chat-service/common/ports"
	"chat-service/internal/types"
	"chat-service/internal/types/actions"
	"context"
	"log"

	"github.com/spf13/viper"
)

type MessageClient struct {
	Hub     *MessageHub
	Conn    ports.ClientTransport
	UserUid string
	ChatUid string
	Send    chan types.Action
}

func NewMessagesClient(
	hub *MessageHub,
	Conn ports.ClientTransport,
	UserUid string, ChatUid string,
) *MessageClient {

	sendChan := make(
		chan types.Action,
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
	msg actions.SendMessageAction,
) {

	message := types.NewMessage(
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
	msg types.Message,
) error {

	data := actions.InitRequestMessageAction(msg)
	c.Send <- types.BuildAction(data)

	return nil
}

func (c *MessageClient) RequestMessageHistory(limit, offset int) {
	history, err := c.Hub.messages.GetMessagesHistory(c.ChatUid, limit, offset)
	if err != nil {
		log.Println("error while fetching history:", err)
		return
	}

	for _, msg := range history {
		data := actions.InitRequestMessageAction(msg)
		c.Send <- types.BuildAction(data)
	}
}

func (c *MessageClient) GetMe() types.User {
	panic("implement me")
}
