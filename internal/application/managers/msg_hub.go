package managers

import (
	"chat-service/internal/application/constants"
	"chat-service/internal/application/ports"
	"chat-service/internal/application/services"
	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/events"
	"chat-service/internal/schema/dto"
	"context"
	"log"
	"slices"

	"sync"
)

type MessagesHub struct {
	broker     events.BrokerMessagesAdaptor
	clients    map[string]map[string]*MessagesClient
	messages   *services.MessageService
	mu         sync.RWMutex
	countUsers int
	msgCount   int
}

func NewMessagesHub(
	messages *services.MessageService,
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

	client.SendMessagesHistory(10, 0)
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

//

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

	sendChan := make(chan dto.SendMessageToClientPayload, 100)

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
	c.Hub.msgCount++
	log.Println("sent to users: ", c.Hub.msgCount)

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
