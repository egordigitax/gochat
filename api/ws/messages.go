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


type MessageHandlerFunc func(
	ctx context.Context,
	data interface{},
	client *messages.MessageClient,
) error

type MessageResponseFunc func(
	ctx context.Context,
	data resources.Action,
	client *messages.MessageClient,
) error

type MessagesWSController struct {
	hub       *messages.MessageHub
	handlers  map[ActionType]MessageHandlerFunc
	responses map[ActionType]MessageResponseFunc
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

	m.handlers = map[ActionType]MessageHandlerFunc{
		ActionType(resources.SEND_MESSAGE): m.HandleSendMessageAction,
	}

	m.responses = map[ActionType]MessageResponseFunc{
		ActionType(resources.REQUEST_MESSAGE): m.ResponseRequestMessageAction,
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
	ctx := context.Background()

	defer func() {
		client.Hub.UnregisterClient(client)
	}()

	for msg := range client.Send {
		responseHandler, ok := m.responses[ActionType(msg.Action)]
		if !ok {
			log.Println("wrong type response recieved")
		}
		responseHandler(ctx, msg, client)

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

func (m *MessagesWSController) ResponseRequestMessageAction(
	ctx context.Context,
	data resources.Action,
	client *messages.MessageClient,
) error {

	actionData, ok := data.Data.(dto.RequestMessagePayload)
	if !ok {
		client.Conn.WriteJSON("Error while handling RequestMessage")
	}

	response := SendMessageToClientResponse{
		Text:      actionData.Text,
		AuthorId:  actionData.AuthorUid,
		Nickname:  "unimplemented",
		CreatedAt: actionData.CreatedAt,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		client.Conn.WriteJSON(fmt.Sprintf("Error while handling RequestMessage: %s", err.Error()))
	}

    rootJson, err := PackToRootMessage(jsonResponse, actionData)

	return client.Conn.WriteMessage(1, rootJson)
}
