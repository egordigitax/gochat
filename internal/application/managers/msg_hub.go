package managers

import (
	"chat-service/internal/application/ports"
	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/events"
	"chat-service/internal/domain/repositories"
	"context"
	"log"
	"sync"
)

type MessagesHub struct {
	broker     events.BrokerMessagesAdaptor
	clients    map[string]map[string]*MessagesClient
	register   chan *MessagesClient
	unregister chan *MessagesClient
	mu         sync.RWMutex
	countUsers int
	msgCount   int
}

func NewMessagesHub(
	repo repositories.MessagesStorage,
	broker events.BrokerMessagesAdaptor,
) *MessagesHub {

	hub := &MessagesHub{
		broker:     broker,
		clients:    make(map[string]map[string]*MessagesClient),
		register:   make(chan *MessagesClient),
		unregister: make(chan *MessagesClient),
	}

	go hub.PumpChats()

	return hub
}

func (h *MessagesHub) PumpChats() {
	ctx, cancel := context.WithCancel(context.Background())

	msgChan, err := h.broker.GetMessagesFromChats(ctx, "chats")
	if err != nil {
		log.Println(err)
	}

	defer func() {
		cancel()
	}()

	for {
		msg := <-msgChan
		h.msgCount++
		log.Println("totalMessages: ", h.msgCount)
		for _, user := range h.clients[msg.ChatUid] {
			select {
			case user.Send <- msg:
			default:
				h.UnregisterClient(user)
			}
		}
	}
}

func (h *MessagesHub) RegisterClient(client *MessagesClient) {
	//TODO: Check if client has access to this chat
	// h.register <- client

	h.countUsers++
	log.Println(h.countUsers)

	client.Send = make(chan entities.Message, 1000)

	go client.WritePump()
	go client.ReadPump()

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.ChatUid] == nil {
		h.clients[client.ChatUid] = make(map[string]*MessagesClient)
	}
	h.clients[client.ChatUid][client.UserUid] = client
}

func (h *MessagesHub) UnregisterClient(client *MessagesClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client.ChatUid][client.UserUid]; ok {
		delete(h.clients[client.ChatUid], client.UserUid)
		close(client.Send)
	}
}

type MessagesClient struct {
	Hub     *MessagesHub
	Conn    ports.ClientTransport
	UserUid string
	ChatUid string
	Send    chan entities.Message
}

func (c *MessagesClient) ReadPump() {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
		cancel()
	}()

	for {
		var msg entities.Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("[ERROR] WebSocket Read:", err)
			break
		}

		msg.UserUid = c.UserUid
		msg.ChatUid = c.ChatUid

		err = c.Hub.broker.SendMessageToChat(ctx, "chats", msg)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (c *MessagesClient) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
