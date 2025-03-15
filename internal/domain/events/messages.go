package events

import (
	"chat-service/internal/domain/entities"
	"context"
)

type BrokerMessagesAdaptor interface {
	GetMessagesFromChannel(
		ctx context.Context,
		chats_uids ...string,
	) (chan entities.Message, error)
	SendMessageToChannel(
		ctx context.Context,
		topic string,
		msg entities.Message,
	) error
	GetMessagesFromQueue(
		ctx context.Context,
		topic string,
	) (chan entities.Message, error)
	SendMessageToQueue(
		ctx context.Context,
		topic string,
		msg entities.Message,
	) error
}
