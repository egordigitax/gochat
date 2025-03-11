package ws_api

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/application/managers"
	"chat-service/internal/schema/resources"
	"log"
	"net/http"
)

//TODO: Use worker pool instead goroutines directly
//TODO: Move it to Controller struct

func ServeChatsWebSocket(hub *managers.ChatsHub, w http.ResponseWriter, r *http.Request) {
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
		Hub:    hub,
		Conn:   conn,
		UserID: userID,
		Send:   make(chan []resources.Chat, 10),
	}

	hub.RegisterClient(client)

	go client.ReadPump()
	go client.WritePump()

	log.Printf("[INFO] Main WS connected: user=%s", userID)
}
