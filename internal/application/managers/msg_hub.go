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

	msgChan, err := h.broker.GetMessagesFromChannel(ctx, constants.CHATS_CHANNEL)
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

	//TODO: Check if client has access to this chat

	h.countUsers++
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.ChatUid] == nil {
		h.clients[client.ChatUid] = make(map[string]*MessagesClient)
	}

	h.clients[client.ChatUid][client.UserUid] = client

	history, err := h.messages.GetMessagesHistory(client.ChatUid, 10, 0)
	if err != nil {
        log.Println("smth wrong:", err)
	}

    // TODO: move to GetMessagesHistory method and handle ASC/DESC orders
    slices.Reverse(history)

	for _, msg := range history {
		client.Send <- dto.BuildSendMessageToClientPayloadFromEntity(msg)
	}
}

func (h *MessagesHub) UnregisterClient(client *MessagesClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client.ChatUid][client.UserUid]; ok {
		delete(h.clients[client.ChatUid], client.UserUid)
		client.Conn.Close()
		close(client.Send)
	}
}

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

	sendChan := make(chan dto.SendMessageToClientPayload, 1000)

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

	// Get Message From Client Logic

	c.Hub.msgCount++
	log.Println("sent to users: ", c.Hub.msgCount)

	message := msg.ToEntity()

	err := c.Hub.broker.SendMessageToChannel(
		ctx,
		constants.CHATS_CHANNEL,
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

	// Send Message To Client Logic

	c.Send <- dto.BuildSendMessageToClientPayloadFromEntity(msg)

	return nil
}
