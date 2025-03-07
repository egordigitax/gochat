package hubs

import (
	"chat-service/internal/domain"
	"chat-service/internal/domain/interfaces"
	"log"
	"sync"
)

type MessagesHub struct {
    //                [ChatUid]   [UserUid]
	clients         map[string]map[string]*MessagesClient
	broadcast       chan domain.Message
	register        chan *MessagesClient
	unregister      chan *MessagesClient
	messagesStorage interfaces.MessagesStorage
	mu              sync.RWMutex
}

func NewMessagesHub(repo interfaces.MessagesStorage) *MessagesHub {
	return &MessagesHub{
		clients:         make(map[string]map[string]*MessagesClient),
		broadcast:       make(chan domain.Message, 100),
		register:        make(chan *MessagesClient),
		unregister:      make(chan *MessagesClient),
		messagesStorage: repo,
	}
}

func (h *MessagesHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
			log.Printf("[INFO] New client joined: user=%s, chat=%s", client.UserUid, client.ChatUid)

		case client := <-h.unregister:
			h.removeClient(client)
			log.Printf("[INFO] Client left: user=%s, chat=%s", client.UserUid, client.ChatUid)

		case message := <-h.broadcast:
			log.Printf("[INFO] Broadcasting message in chat=%s: %s", message.ChatUid, message.Text)
			h.sendMessage(message)
            go h.messagesStorage.SaveMessage(message)
		}
	}
}

func (h *MessagesHub) RegisterClient(client *MessagesClient) {
    //TODO: Check if client has access to this chat
	h.register <- client
}

func (h *MessagesHub) UnregisterClient(client *MessagesClient) {
	h.unregister <- client
}

func (h *MessagesHub) addClient(client *MessagesClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[client.ChatUid] == nil {
		h.clients[client.ChatUid] = make(map[string]*MessagesClient)
	}
	h.clients[client.ChatUid][client.UserUid] = client
}

func (h *MessagesHub) removeClient(client *MessagesClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client.ChatUid][client.UserUid]; ok {
		delete(h.clients[client.ChatUid], client.UserUid)
		close(client.Send)
	}
}

func (h *MessagesHub) sendMessage(message domain.Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.clients[message.ChatUid] {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients[message.ChatUid], client.UserUid)
		}
	}
}

type MessagesClient struct {
	Hub     *MessagesHub
	Conn    interfaces.ClientTransport
	UserUid string
	ChatUid string
	Send    chan domain.Message
}

func (c *MessagesClient) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg domain.Message
		err := c.Conn.ReadJSON(&msg)

		if err != nil {
			log.Println("[ERROR] WebSocket Read:", err)
			break
		}

		log.Printf("[INFO] Received message from %s in chat %s: %s", c.UserUid, c.ChatUid, msg.Text)

		msg.UserUid = c.UserUid
		msg.ChatUid = c.ChatUid

		c.Hub.broadcast <- msg
	}
}

func (c *MessagesClient) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
