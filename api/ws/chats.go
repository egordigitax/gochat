package ws_api

import (
	"chat-service/internal/application/schema/dto"
	"chat-service/internal/application/use_cases/chat_list"
	"chat-service/internal/utils"
	"log"
	"net/http"
)

//TODO: Use worker pool instead goroutines directly

type ChatsWSController struct {
	hub *chat_list.ChatsHub
}

func NewChatsWSController(
	hub *chat_list.ChatsHub,
) *ChatsWSController {
	return &ChatsWSController{
		hub: hub,
	}
}

func (c *ChatsWSController) Handle() {
	http.HandleFunc("/chats", func(w http.ResponseWriter, r *http.Request) {
		c.ServeChatsWebSocket(w, r)
	})
}

func (c *ChatsWSController) ServeChatsWebSocket(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	upgrader := GetUpgrader()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] WebSocket upgrade failed:", err)
		return
	}

	client := chat_list.NewChatsClient(c.hub, conn, userID)

	c.hub.RegisterClient(client)

	go c.StartClientWrite(client)
}

func (c *ChatsWSController) StartClientWrite(client *chat_list.ChatsClient) {
	defer func() {
		c.hub.UnregisterClient(client)
	}()

	for msg := range client.Send {
		// handle different actions and parse to schema
		message, ok := msg.Data.(dto.GetUserChatsPayload)
		if !ok {
			log.Println("[ERROR] wrong type of data")
			continue
		}

		items := make([]Chat, len(message.Items))
		for i, item := range message.Items {
			items[i] = Chat{
				Title:       item.Title,
				UnreadCount: item.UnreadCount,
				LastMessage: item.LastMessage.Text,
				LastAuthor:  item.LastMessage.AuthorUid,
				MediaUrl:    item.MediaUrl,
			}
		}

		response := GetChatsResponse{
			ActionType: "get_chats",
			Items:      items,
		}

		if err := client.Conn.WriteJSON(response); err != nil {
			break
		}
	}
}
