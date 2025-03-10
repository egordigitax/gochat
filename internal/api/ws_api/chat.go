package ws_api

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/application/managers"
	"log"
	"net/http"
)

//TODO: Use worker pool instead goroutines directly
//TODO: Move it to Controller struct

func ServeWebSocket(hub *managers.MessagesHub, w http.ResponseWriter, r *http.Request) {
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

	upgrader := utils.GetUpgrader()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket error:", err)
		return
	}

	client := &managers.MessagesClient{
		Hub:     hub,
		Conn:    conn,
		UserUid: userID,
		ChatUid: chatID,
	}

	hub.RegisterClient(client)
    log.Println("user ", userID, " connected.")
}
