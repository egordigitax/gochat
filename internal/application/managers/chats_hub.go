package managers

import (
	"chat-service/internal/application/constants"
	"chat-service/internal/application/ports"
	"chat-service/internal/application/services"
	"chat-service/internal/domain/events"
	"chat-service/internal/schema/dto"
	"context"
	"log"
	"sync"
	"time"
)

type ChatsHub struct {
	broker          events.BrokerMessagesAdaptor
	clients         map[string]*ChatsClient
	isChatHasClient map[string]map[string]bool
	messages        *services.MessageService
	chats           *services.ChatsService
	mu              sync.RWMutex
}

func NewChatsHub(
	messagesService *services.MessageService,
	chatsService *services.ChatsService,
	messagesBroker events.BrokerMessagesAdaptor,
) *ChatsHub {
	return &ChatsHub{
		clients:  make(map[string]*ChatsClient),
		messages: messagesService,
		chats:    chatsService,
		broker:   messagesBroker,
	}
}

func (h *ChatsHub) StartPumpChats() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	msgChan, err := h.broker.GetMessagesFromChannel(ctx, constants.SAVED_MESSAGES_CHANNEL)
	if err != nil {
		log.Println("Cant subscribe to chats:", err)
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	chats := make(map[string]struct{})

	for {

		select {
		case msg := <-msgChan:
			chats[msg.ChatUid] = struct{}{}

		case <-ticker.C:
			for chat := range chats {
				userUids, err := h.chats.GetAllUsersFromChatByUid(chat)
				if err != nil {
					log.Println(err)
					continue
				}
				for _, user := range userUids {
					if client, ok := h.clients[user]; ok {
						client.UpdateChats()
					}
				}
			}
		}
	}
}

func (h *ChatsHub) RegisterClient(client *ChatsClient) {
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

	h.clients[client.UserID].Send <- chats

}

func (h *ChatsHub) UnregisterClient(client *ChatsClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client.UserID)
	close(client.Send)
}

type ChatsClient struct {
	Hub    *ChatsHub
	Conn   ports.ClientTransport
	UserID string
	Send   chan dto.GetUserChatsByUidResponse
}

func (c *ChatsClient) UpdateChats() {
	chats, err := c.Hub.chats.GetChatsByUserUid(
		dto.GetUserChatsByUidPayload{
			UserUid: c.UserID,
		})
	if err != nil {
		log.Println(err)
		return
	}
	c.Send <- chats
}

func (c *ChatsClient) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
