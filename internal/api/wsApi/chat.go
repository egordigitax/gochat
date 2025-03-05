package wsApi

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/domain"
	"chat-service/internal/infra"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (for testing)
	}}

func ServeWebSocket(hub *infra.Hub, w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromHeader(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		http.Error(w, "Missing chat_id", http.StatusBadRequest)
		return
	} else {
		fmt.Println("got message to chat_id: ", r.URL.Query().Get("chat_id"))
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket error:", err)
		return
	}

	client := &infra.Client{Hub: hub, Conn: conn, UserID: userID, ChatID: chatID, Send: make(chan domain.Message, 10)}
	hub.RegisterClient(client)

	go client.ReadPump()  // This listens for messages from the client
	go client.WritePump() // This sends messages to the client
}
