package services

import (
	"chat-service/internal/domain"
	"chat-service/internal/domain/interfaces"
	"log"
	"sync"
)

type Hub struct {
	clients         map[string]map[string]*Client
	broadcast       chan domain.Message
	register        chan *Client
	unregister      chan *Client
	messagesStorage interfaces.MessagesStorage
	mu              sync.RWMutex
}

func NewHub(repo interfaces.MessagesStorage) *Hub {
	return &Hub{
		clients:         make(map[string]map[string]*Client),
		broadcast:       make(chan domain.Message, 100),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		messagesStorage: repo,
	}
}

func (h *Hub) Run() {
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

func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *Hub) addClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[client.ChatUid] == nil {
		h.clients[client.ChatUid] = make(map[string]*Client)
	}
	h.clients[client.ChatUid][client.UserUid] = client
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client.ChatUid][client.UserUid]; ok {
		delete(h.clients[client.ChatUid], client.UserUid)
		close(client.Send)
	}
}

func (h *Hub) sendMessage(message domain.Message) {
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

type Client struct {
	Hub     *Hub
	Conn    interfaces.ClientTransport
	UserUid string
	ChatUid string
	Send    chan domain.Message
}

func (c *Client) ReadPump() {
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

func (c *Client) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
