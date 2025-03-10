package managers

import (
	"chat-service/internal/application/ports"
	"chat-service/internal/domain/entities"
	"chat-service/internal/domain/events"
	"chat-service/internal/domain/repositories"
	"log"
	"sync"
)

type MessagesHub struct {
	broker     events.BrokerMessagesAdaptor
	clients    map[string]map[string]*MessagesClient
	register   chan *MessagesClient
	unregister chan *MessagesClient
	mu         sync.RWMutex
}

func NewMessagesHub(
	repo repositories.MessagesStorage,
	broker events.BrokerMessagesAdaptor,
) *MessagesHub {
	return &MessagesHub{
		broker:     broker,
		clients:    make(map[string]map[string]*MessagesClient),
		register:   make(chan *MessagesClient),
		unregister: make(chan *MessagesClient),
	}
}

func (h *MessagesHub) RegisterClient(client *MessagesClient) {
	//TODO: Check if client has access to this chat
	msgChan, err := h.broker.GetMessagesFromChats(client.ChatUid)
	if err != nil {
		log.Println(err)
	}

	client.Send = msgChan

	// h.register <- client

    go client.WritePump()
    go client.ReadPump()
}

func (h *MessagesHub) UnregisterClient(client *MessagesClient) {
	// h.unregister <- client
}

type MessagesClient struct {
	Hub     *MessagesHub
	Conn    ports.ClientTransport
	UserUid string
	ChatUid string
	Send    chan entities.Message
}

func (c *MessagesClient) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
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

		err = c.Hub.broker.SendMessageToChat(msg)
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
