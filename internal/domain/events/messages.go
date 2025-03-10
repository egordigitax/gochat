package events

import "chat-service/internal/domain/entities"

type BrokerMessagesAdaptor interface {
	GetMessagesFromChats(
		chats_uids ...string,
	) (chan entities.Message, error)
	SendMessageToChat(
		msg entities.Message,
	) error
}
