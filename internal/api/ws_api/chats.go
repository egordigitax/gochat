package ws_api

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/application/managers"
	"chat-service/internal/schema/dto"
	"log"
	"net/http"
)

//TODO: Use worker pool instead goroutines directly

type ChatsWSController struct {
	hub *managers.ChatsHub
}

func NewChatsWSController(
	hub *managers.ChatsHub,
) *ChatsWSController {
	return &ChatsWSController{
		hub: hub,
	}
}

func (c *ChatsWSController) ServeChatsWebSocket(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	upgrader := utils.GetUpgrader()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] WebSocket upgrade failed:", err)
		return
	}

	client := &managers.ChatsClient{
		Hub:    c.hub,
		Conn:   conn,
		UserID: userID,
		Send:   make(chan dto.GetUserChatsByUidResponse, 10),
	}

	c.hub.RegisterClient(client)

	// go c.StartClientRead(client)
	go c.StartClientWrite(client)
}

func (c *ChatsWSController) StartClientWrite(client *managers.ChatsClient) {
	defer func() {
		c.hub.UnregisterClient(client)
	}()

	for msg := range client.Send {
		if err := client.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
//
// func (c *ChatsWSController) StartClientRead(client *managers.ChatsClient) {
// 	defer func() {
// 		c.hub.UnregisterClient(client)
// 	}()
//
// }
