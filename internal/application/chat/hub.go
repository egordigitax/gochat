package chat

import (
	"chat-service/internal/application/common/constants"
	"chat-service/internal/domain/events"
	"chat-service/internal/schema/dto"
	"context"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type ChatsHub struct {
	broker          events.BrokerMessagesAdaptor
	clients         map[string]*ChatsClient
	isChatHasClient map[string]map[string]bool
	chats           IChatsService
	mu              sync.RWMutex
}

func NewChatsHub(
	chatsService IChatsService,
	messagesBroker events.BrokerMessagesAdaptor,
) *ChatsHub {
	return &ChatsHub{
		clients: make(map[string]*ChatsClient),
		chats:   chatsService,
		broker:  messagesBroker,
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

	ticker := time.NewTicker(viper.GetDuration("app.chats_update_rate") * time.Millisecond)
	chats := make(map[string]struct{})

	log.Println("Chats pump started")

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
				delete(chats, chat)
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
