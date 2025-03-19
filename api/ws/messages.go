package ws_api

import (
	"chat-service/internal/application/schema/dto"
	"chat-service/internal/application/schema/resources"
	"chat-service/internal/application/use_cases/messages"
	"chat-service/internal/utils"
	"context"
	"encoding/json"
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
        log.Println(msg.Data)
		if msg.Action == resources.REQUEST_MESSAGE {
			actionData, _ := msg.Data.(dto.RequestMessagePayload)
            
            log.Println(actionData.Text)
			data := SendMessageToClientResponse{
				Text:      actionData.Text,
				AuthorId:  actionData.AuthorUid,
				Nickname:  "implement me",
				CreatedAt: actionData.CreatedAt,
			}

			payload, err := json.Marshal(data)
			if err != nil {
				log.Println(err)
			}

			response := RootMessage{
				ActionType: ActionType(resources.REQUEST_MESSAGE),
				RawPayload: payload,
			}

			if err := client.Conn.WriteJSON(response); err != nil {
				break
			}
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
		var root RootMessage
		err := client.Conn.ReadJSON(&root)
		if err != nil {
			err = client.Conn.WriteJSON(fmt.Sprintf("Error: %s", err))
			if err != nil {
				break
			}
		}

		switch root.ActionType {
		case ActionType(resources.SEND_MESSAGE):
			var msg GetMessageFromClientRequest
			if err := json.Unmarshal(root.RawPayload, &msg); err != nil {
				log.Println("Failed to unmarshal GetMessageFromClientRequest:", err)
				return
			}
			// Handle msg
			client.SendMessage(ctx, dto.SendMessagePayload{
				ChatUid:   client.ChatUid,
				AuthorUid: client.UserUid,
				CreatedAt: time.Now().String(),
				Text:      msg.Text,
			})

		case "send_message_to_client":
			var msg SendMessageToClientResponse
			if err := json.Unmarshal(root.RawPayload, &msg); err != nil {
				log.Println("Failed to unmarshal SendMessageToClientResponse:", err)
				return
			}
			// Handle msg

		case "get_chats":
			var msg GetChatsResponse
			if err := json.Unmarshal(root.RawPayload, &msg); err != nil {
				log.Println("Failed to unmarshal GetChatsResponse:", err)
				return
			}
			// Handle msg

		default:
			log.Println("Unknown action type:", root.ActionType)
		}

		// TODO: Handle different action types from user
	}
}
