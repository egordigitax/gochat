package managers

import (
	"chat-service/internal/application/ports"
	"chat-service/internal/application/services"
	"chat-service/internal/domain/events"
	"chat-service/internal/schema/dto"
	"context"
	"log"
	"sync"
)

type ChatsHub struct {
	broker       events.BrokerMessagesAdaptor
	clients      map[string]*ChatsClient
	clientsChats map[string][]string
	messages     *services.MessageService
	chats        *services.ChatsService
	mu           sync.RWMutex
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

func (h *ChatsHub) Run() {

	ctx, cancel := context.WithCancel(context.Background())

	msgChan, err := h.broker.GetMessagesFromChats(ctx, "chats")
	if err != nil {
		log.Println("Cant subscribe to chats: ", err.Error())
	}

	defer func() {
		cancel()
		close(msgChan)
	}()

	for {
		select {
		case msg := <-msgChan:
			if client, ok := h.clients[msg.ChatUid]; ok {
				client.UpdateChats()
			}
		default:
			break
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

	// TODO: map by chatUid instead of useruid

	h.clientsChats[client.UserID] = make([]string, len(chats.Items))

	for i, item := range chats.Items {
		h.clientsChats[client.UserID][i] = item.Uid
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
