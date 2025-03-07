package ws_api

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/application/hubs"
	"chat-service/internal/domain"
	"log"
	"net/http"
)

//TODO: Use worker pool instead goroutines directly
//TODO: Move it to Controller struct

func ServeMainWebSocket(hub *hubs.ChatsHub, w http.ResponseWriter, r *http.Request) {
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

	client := &hubs.ChatsClient{
		Hub:    hub,
		Conn:   conn,
		UserID: userID,
		Send:   make(chan []domain.Chat, 10),
	}

	client.Done = make(chan struct{}, 1)
	hub.RegisterClient(client)

	go client.ReadPump()
	go client.WritePump()

	log.Printf("[INFO] Main WS connected: user=%s", userID)
}
