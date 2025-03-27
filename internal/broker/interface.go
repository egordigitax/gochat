package broker

import (
	"chat-service/internal/types"
	"context"
)

type BaseAdaptor interface {
	Subscribe(ctx context.Context, topics ...string) (chan string, error)
	Publish(ctx context.Context, topic, message string) error
	ToQueue(ctx context.Context, topic, message string) error
	FromQueue(ctx context.Context, topic string) (string, error)
}

type MessagesAdaptor interface {
	GetMessagesFromChannel(
		ctx context.Context,
		chats_uids ...string,
	) (chan types.Message, error)
	SendMessageToChannel(
		ctx context.Context,
		topic string,
		msg types.Message,
	) error
	GetMessagesFromQueue(
		ctx context.Context,
		topic string,
	) (chan types.Message, error)
	SendMessageToQueue(
		ctx context.Context,
		topic string,
		msg types.Message,
	) error
}
