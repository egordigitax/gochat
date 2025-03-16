package message

import (
	"chat-service/internal/application/constants"
	"chat-service/internal/domain/events"
	"context"
	"log"

	"sync"
)

type MessagesHub struct {
	broker     events.BrokerMessagesAdaptor
	clients    map[string]map[string]*MessagesClient
	messages   IMessagesService
	mu         sync.RWMutex
	countUsers int
	msgCount   int
}

func NewMessagesHub(
	messages IMessagesService,
	broker events.BrokerMessagesAdaptor,
) *MessagesHub {

	hub := &MessagesHub{
		broker:   broker,
		messages: messages,
		clients:  make(map[string]map[string]*MessagesClient),
	}

	return hub
}

func (h *MessagesHub) StartPumpMessages() {
	ctx, cancel := context.WithCancel(context.Background())

	msgChan, err := h.broker.GetMessagesFromChannel(ctx, constants.SAVED_MESSAGES_CHANNEL)
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
			user.SendMessageToClient(ctx, msg)
		}
		h.mu.RUnlock()
	}
}

func (h *MessagesHub) RegisterClient(client *MessagesClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// h.checkIfUserHasPrevConnectionUnsafe(client)

	if h.clients[client.ChatUid] == nil {
		h.clients[client.ChatUid] = make(map[string]*MessagesClient)
	}

	h.clients[client.ChatUid][client.UserUid] = client

	client.SendMessagesHistory(3, 0)
}

func (h *MessagesHub) UnregisterClient(client *MessagesClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.unregisterClientUnsafe(client)
}

// Unsafe methonds should be called only with Mutex Lock

func (h *MessagesHub) isUserExistUnsafe(client *MessagesClient) (*MessagesClient, bool) {
	client, ok := h.clients[client.ChatUid][client.UserUid]
	return client, ok
}

func (h *MessagesHub) checkIfUserHasPrevConnectionUnsafe(client *MessagesClient) {
	if oldClient, ok := h.isUserExistUnsafe(client); ok {
		h.unregisterClientUnsafe(oldClient)
	}
}

func (h *MessagesHub) unregisterClientUnsafe(client *MessagesClient) {
	oldClient, ok := h.isUserExistUnsafe(client)
	if !ok || oldClient != client {
		return
	}

	delete(h.clients[client.ChatUid], client.UserUid)
	oldClient.Conn.Close()
	close(oldClient.Send)

	if len(h.clients[client.ChatUid]) == 0 {
		delete(h.clients, client.ChatUid)
	}
}
