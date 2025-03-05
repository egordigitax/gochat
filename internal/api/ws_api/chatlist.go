package ws_api

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/application/services"
	"chat-service/internal/domain"
	"log"
	"net/http"
)

func ServeMainWebSocket(hub *services.ChatsHub, w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] WebSocket upgrade failed:", err)
		return
	}

	client := &services.ChatsClient{
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
