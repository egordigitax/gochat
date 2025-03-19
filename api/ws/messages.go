package ws_api

import (
	"chat-service/internal/application/schema/dto"
	"chat-service/internal/application/use_cases/messages"
	"chat-service/internal/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// TODO: Use worker pool instead goroutines directly
// TODO: Move it to Controller struct

type MessagesWSController struct {
	hub *messages.MessageHub
}

func NewMessagesWSController(
	hub *messages.MessageHub,
) *MessagesWSController {
	return &MessagesWSController{
		hub: hub,
	}
}

func (m *MessagesWSController) Handle() {
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		m.ServeMessagesWebSocket(w, r)
	})
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

	upgrader := GetUpgrader()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket error:", err)
		return
	}

	client := messages.NewMessagesClient(
		m.hub,
		conn,
		userID,
		chatID,
	)

	m.hub.RegisterClient(client)

	go m.StartClientWrite(client)
	go m.StartClientRead(client)
}

func (m *MessagesWSController) StartClientWrite(
	client *messages.MessageClient,
) {

	defer func() {
		client.Hub.UnregisterClient(client)
	}()

	for msg := range client.Send {

		// handle different actions and parse to schema
		data, _ := msg.Data.(dto.GetMessagePayload)

		message := SendMessageToClientResponse{
			ActionType: ActionType("get_message"),
			Text:       data.Text,
			AuthorId:   data.AuthorUid,
			Nickname:   "implement me",
			CreatedAt:  data.CreatedAt,
		}

		if err := client.Conn.WriteJSON(message); err != nil {
			break
		}
	}
}

func (m *MessagesWSController) StartClientRead(
	client *messages.MessageClient,
) {

	// TODO: test cancel, and add it to defer if it works fine

	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		client.Hub.UnregisterClient(client)
		cancel()
	}()

	for {

		// handle different actions and parse to schema
		var msg GetMessageFromClientRequest
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			err = client.Conn.WriteJSON(fmt.Sprintf("Error: %s", err))
			if err != nil {
				break
			}
		}

		// TODO: Handle different action types from user

		client.SendMessage(ctx, dto.SendMessagePayload{
			ChatUid:   client.ChatUid,
			AuthorUid: client.UserUid,
			CreatedAt: time.Now().String(),
			Text:      msg.Text,
		})
	}
}
