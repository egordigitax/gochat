package managers

import (
	"chat-service/internal/application/ports"
	"chat-service/internal/application/services"
	"chat-service/internal/domain/entities"
	"chat-service/internal/schema/dto"
	"chat-service/internal/schema/resources"
	"log"
	"sync"
	"time"
)

type ChatsHub struct {
	clients    map[string]*ChatsClient
	broadcast  chan entities.Message
	register   chan *ChatsClient
	unregister chan *ChatsClient
	messages   *services.MessageService
	chats      *services.ChatsService
	mu         sync.RWMutex
}

func NewChatsHub(
	messagesService *services.MessageService,
	chatsService *services.ChatsService,
) *ChatsHub {
	return &ChatsHub{
		clients:    make(map[string]*ChatsClient),
		broadcast:  make(chan entities.Message, 1),
		register:   make(chan *ChatsClient),
		unregister: make(chan *ChatsClient),
		messages:   messagesService,
		chats:      chatsService,
	}
}

func (h *ChatsHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

			go func(c *ChatsClient) {
				for {
					select {
					case <-c.Done:
						return
					default:
						chats, err := h.chats.GetChatsByUserUid(
							dto.GetUserChatsByUidPayload{
								UserUid: client.UserID,
							})
						if err != nil {
							log.Println(err)
							return
						}
						h.clients[client.UserID].Send <- chats.Items
						time.Sleep(1 * time.Second)

					}
				}
			}(client)

		case client := <-h.unregister:
			h.mu.Lock()
			delete(h.clients, client.UserID)
			close(client.Send)
			close(client.Done)
			h.mu.Unlock()
		}
	}
}

func (h *ChatsHub) RegisterClient(client *ChatsClient) {
	h.register <- client
}

func (h *ChatsHub) UnregisterClient(client *ChatsClient) {
	client.Done <- struct{}{}
	h.unregister <- client
}

type ChatsClient struct {
	Hub    *ChatsHub
	Conn   ports.ClientTransport // move from repos
	UserID string
	Send   chan []resources.Chat
	Done   chan struct{}
}

func (c *ChatsClient) ReadPump() {
	defer func() {
		c.Hub.UnregisterClient(c)
		c.Conn.Close()
	}()

	for {
		var msg entities.Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("[ERROR] WebSocket Read:", err)
			break
		}

		c.Hub.broadcast <- msg
	}

	log.Println("[DEBUG] Main ReadPump started for user:", c.UserID)
}

func (c *ChatsClient) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
