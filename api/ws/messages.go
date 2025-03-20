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
	hub      *messages.MessageHub
	handlers map[ActionType]HandlerFunc
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

	m.handlers = map[ActionType]HandlerFunc{
		ActionType(resources.SEND_MESSAGE): m.HandleSendMessageAction,
	}

    // TODO: add responses instead of standalone serializer (?)
}

func (m *MessagesWSController) ServeMessagesWebSocket(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIDFromHeader(
		r.Header.Get("Authorization"),
	)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chatId := r.URL.Query().Get("chat_id")
	if chatId == "" {
		http.Error(w, "Missing chat_id", http.StatusBadRequest)
		return
	}

	upgrader := utils.GetUpgrader()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket error:", err)
		return
	}

	client := messages.NewMessagesClient(
		m.hub,
		conn,
		userId,
		chatId,
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
		response, err := Serialize(msg)
		if err != nil {
			log.Println(err)
		}
		err = client.Conn.WriteJSON(response)
		if err != nil {
			break
		}
	}
}

func (m *MessagesWSController) StartClientRead(
	client *messages.MessageClient,
) {
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		client.Hub.UnregisterClient(client)
		cancel()
	}()

	for {
		var root RootMessage
		err := client.Conn.ReadJSON(&root)
		if err != nil {
			err = client.Conn.WriteJSON(fmt.Sprintf("Error: %s", err))
			if err != nil {
				break
			}
		}

		handler, ok := m.handlers[root.ActionType]
		if !ok {
			log.Println("unknown action type")
		}

		err = handler(ctx, root.RawPayload, client)
		if err != nil {
			log.Println("error while handling: ", root.ActionType)
		}
	}
}

func (m *MessagesWSController) HandleSendMessageAction(
	ctx context.Context,
	data interface{},
	client *messages.MessageClient,
) error {
	var msg GetMessageFromClientRequest
	if err := json.Unmarshal(data.(json.RawMessage), &msg); err != nil {
		log.Println("Failed to unmarshal GetMessageFromClientRequest:", err)
		return err
	}

	client.SendMessage(ctx, dto.SendMessagePayload{
		ChatUid:   client.ChatUid,
		AuthorUid: client.UserUid,
		CreatedAt: time.Now().String(),
		Text:      msg.Text,
	})

	return nil
}
