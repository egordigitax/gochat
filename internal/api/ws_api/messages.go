package ws_api

import (
	"chat-service/internal/api/utils"
	"chat-service/internal/application/managers"
	"chat-service/internal/schema/dto"
	"context"
	"fmt"
	"log"
	"net/http"
)

// TODO: Use worker pool instead goroutines directly
// TODO: Move it to Controller struct
type MessagesWSController struct {
	hub *managers.MessagesHub
}

func NewMessagesWSController(
	hub *managers.MessagesHub,
) *MessagesWSController {
	return &MessagesWSController{
		hub: hub,
	}
}

func (m *MessagesWSController) ServeMessagesWebSocket(w http.ResponseWriter, r *http.Request) {
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
		Hub:     m.hub,
		Conn:    conn,
		UserUid: userID,
		ChatUid: chatID,
	}

	m.hub.RegisterClient(client)

	go m.StartClientWrite(client)
	go m.StartClientRead(client)
}

func (m *MessagesWSController) StartClientWrite(
	client *managers.MessagesClient,
) {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		client.Hub.UnregisterClient(client) // move somewhere
		cancel()
	}()

	for msg := range client.Send {

		response := client.GetMessageFromChat(ctx, msg)

		if err := client.Conn.WriteJSON(response); err != nil {
			break
		}
	}
}

func (m *MessagesWSController) StartClientRead(
	client *managers.MessagesClient,
) {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		client.Hub.UnregisterClient(client)
		cancel()
	}()

	for {
		var msg dto.SendMessageToChatPayload
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println("[ERROR] WebSocket Read:", err)
			break
		}

		msg.Msg.AuthorUid = client.UserUid
		msg.Msg.ChatUid = client.ChatUid

		err = msg.Validate()
		if err != nil {
			client.Conn.WriteJSON(fmt.Sprintf("Error: %s", err.Error()))
			continue
		}

		client.SendMessageToChat(ctx, msg)
	}
}
