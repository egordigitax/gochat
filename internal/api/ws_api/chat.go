package ws_api

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/application/services"
	"chat-service/internal/domain"
	"log"
	"net/http"
)

func ServeWebSocket(hub *services.Hub, w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromHeader(
		r.Header.Get("Authorization"),
	)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		http.Error(w, "Missing chat_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket error:", err)
		return
	}

	client := &services.Client{
		Hub:    hub,
		Conn:   conn,
		UserID: userID,
		ChatID: chatID,
		Send:   make(chan domain.Message, 10),
	}

	hub.RegisterClient(client)

	go client.ReadPump()
	go client.WritePump()
}
