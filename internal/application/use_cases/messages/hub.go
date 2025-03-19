package messages

import (
	"chat-service/internal/application/common/constants"
	"chat-service/internal/application/schema/resources"
	"chat-service/internal/domain/events"
	"context"
	"log"

	"sync"
)

type MessageHub struct {
	broker   events.BrokerMessagesAdaptor
	clients  map[string]map[string]*MessageClient
	messages IMessageService
	mu       sync.RWMutex
}

func NewMessagesHub(
	messages IMessageService,
	broker events.BrokerMessagesAdaptor,
) *MessageHub {

	hub := &MessageHub{
		broker:   broker,
		messages: messages,
		clients:  make(map[string]map[string]*MessageClient),
	}

	return hub
}

func (h *MessageHub) StartPumpMessages() {
	ctx, cancel := context.WithCancel(context.Background())

	msgChan, err := h.broker.GetMessagesFromChannel(ctx, constants.CHATS_CHANNEL)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		cancel()
	}()

	log.Println("Messages pump started")

	for {
		msg := <-msgChan
		h.mu.RLock()
		clients := h.clients[msg.ChatUid]
		for _, user := range clients {
			err := user.RequestMessage(resources.Message{
				Username:  msg.UserInfo.Nickname,
				AuthorUid: msg.UserUid,
				ChatUid:   msg.ChatUid,
				Text:      msg.Text,
				CreatedAt: msg.CreatedAt,
			})
			if err != nil {
				log.Println(err)
			}
		}
		h.mu.RUnlock()
	}
}

func (h *MessageHub) RegisterClient(client *MessageClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// h.checkIfUserHasPrevConnectionUnsafe(client)

	if h.clients[client.ChatUid] == nil {
		h.clients[client.ChatUid] = make(map[string]*MessageClient)
	}

	h.clients[client.ChatUid][client.UserUid] = client

	client.RequestMessageHistory(10, 0)
}

func (h *MessageHub) UnregisterClient(client *MessageClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.unregisterClientUnsafe(client)
}

// Unsafe methonds should be called only with Mutex Lock

func (h *MessageHub) isUserExistUnsafe(client *MessageClient) (*MessageClient, bool) {
	client, ok := h.clients[client.ChatUid][client.UserUid]
	return client, ok
}

func (h *MessageHub) checkIfUserHasPrevConnectionUnsafe(client *MessageClient) {
	if oldClient, ok := h.isUserExistUnsafe(client); ok {
		h.unregisterClientUnsafe(oldClient)
	}
}

func (h *MessageHub) unregisterClientUnsafe(client *MessageClient) {
	oldClient, ok := h.isUserExistUnsafe(client)
	if !ok || oldClient != client {
		return
	}

	delete(h.clients[client.ChatUid], client.UserUid)
	err := oldClient.Conn.Close()
	if err != nil {
		return
	}
	close(oldClient.Send)

	if len(h.clients[client.ChatUid]) == 0 {
		delete(h.clients, client.ChatUid)
	}
}
