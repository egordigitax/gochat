package infra

import (
	"chat-service/internal/domain"
	"log"
	"sync"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[string]map[string]*Client
	broadcast  chan domain.Message
	register   chan *Client
	unregister chan *Client
	repo       domain.ChatRepository
	mu         sync.RWMutex
}

func NewHub(repo domain.ChatRepository) *Hub {
	return &Hub{
		clients:    make(map[string]map[string]*Client),
		broadcast:  make(chan domain.Message, 100),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		repo:       repo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
			log.Printf("[INFO] New client joined: user=%s, chat=%s", client.UserID, client.ChatID)

		case client := <-h.unregister:
			h.removeClient(client)
			log.Printf("[INFO] Client left: user=%s, chat=%s", client.UserID, client.ChatID)

		case message := <-h.broadcast:
			log.Printf("[INFO] Broadcasting message in chat=%s: %s", message.ChatID, message.Text)
			h.sendMessage(message)
			go h.repo.SaveMessage(message)
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
	if h.clients[client.ChatID] == nil {
		h.clients[client.ChatID] = make(map[string]*Client)
	}
	h.clients[client.ChatID][client.UserID] = client
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client.ChatID][client.UserID]; ok {
		delete(h.clients[client.ChatID], client.UserID)
		close(client.Send)
	}
}

func (h *Hub) sendMessage(message domain.Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.clients[message.ChatID] {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients[message.ChatID], client.UserID)
		}
	}
}

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	UserID string
	ChatID string
	Send   chan domain.Message
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

		log.Printf("[INFO] Received message from %s in chat %s: %s", c.UserID, c.ChatID, msg.Text)

		msg.UserID = c.UserID
		msg.ChatID = c.ChatID

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
