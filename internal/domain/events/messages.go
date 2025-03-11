package events

import (
	"chat-service/internal/domain/entities"
	"context"
)

type BrokerMessagesAdaptor interface {
	GetMessagesFromChats(
		ctx context.Context,
		chats_uids ...string,
	) (chan entities.Message, error)
	SendMessageToChat(
		ctx context.Context,
		topic string,
		msg entities.Message,
	) error
}
