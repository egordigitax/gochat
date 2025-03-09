package managers

import (
	"chat-service/internal/application/ports"
	"chat-service/internal/application/services"
	"chat-service/internal/domain/entities"
	"chat-service/internal/schema/dto"
	"chat-service/internal/schema/resources"
	"log"
	"sync"
)

type ChatsHub struct {
	clients    map[string]*ChatsClient
	broadcast  chan string
	register   chan *ChatsClient
	unregister chan *ChatsClient
	messages   *services.MessageService
	chats      *services.ChatsService
	mu         sync.RWMutex
}

func NewChatsHub(
	messagesService *services.MessageService,
	chatsService *services.ChatsService,
	updateChan *chan string,
) *ChatsHub {
	return &ChatsHub{
		clients:    make(map[string]*ChatsClient),
		broadcast:  *updateChan,
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

			chats, err := h.chats.GetChatsByUserUid(
				dto.GetUserChatsByUidPayload{
					UserUid: client.UserID,
				})
			if err != nil {
				log.Println(err)
				return
			}

			h.clients[client.UserID].Send <- chats.Items

		case chatUid := <-h.broadcast:

			log.Println("got a new message! updating chat for members: ", chatUid)

			userUids, err := h.chats.GetAllUsersFromChatByUid(chatUid)
			if err != nil {
				log.Println("error: ", err)
			}

			log.Println("got users: ", userUids)

			for _, uid := range userUids {
				chats, err := h.chats.GetChatsByUserUid(
					dto.GetUserChatsByUidPayload{
						UserUid: uid,
					})
				if err != nil {
					log.Println(err)
					return
				}
				log.Println("get chats for user :", uid, " chats: ", chats)

				if client, ok := h.clients[uid]; ok {
					client.Send <- chats.Items
				}
			}

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

func (h *ChatsHub) UpdateChatForUsers(chat_uid string) {
	h.broadcast <- chat_uid
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
	}
}

func (c *ChatsClient) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
